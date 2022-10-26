#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

mod commands;
mod database;
use sqlx::sqlite::SqlitePool;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use tauri::{
    api::process::{Command, CommandEvent},
    Manager,
};

#[tauri::command]
async fn send_message(function: String, args: String) -> String {
    let mut map = HashMap::new();
    map.insert("function", function);
    map.insert("args", args);
    let client = reqwest::Client::new();
    let resp = client
        .post("http://localhost:3451/ipc")
        .json(&map)
        .send()
        .await
        .unwrap()
        .text()
        .await
        .unwrap();
    resp
}

fn main() {
    use tauri::async_runtime::block_on;
    const DATABASE_DIR: &str = "/home/dlangevi/.config/read-chinese";
    const DATABASE_NAME: &str = "db.sqlite3";

    let database_url = format!("sqlite:{}/{}", DATABASE_DIR, DATABASE_NAME);
    let pool: SqlitePool = block_on(database::create_pool(&database_url)).expect("cannot connect");

    block_on(database::migrate(&pool)).expect("failed to migrate db");

    tauri::Builder::default()
        .manage(pool)
        .invoke_handler(tauri::generate_handler![
            send_message,
            commands::book_library::learning_target,
            commands::book_library::learning_target_book
        ])
        .setup(|app| {
            let window = app.get_window("main").unwrap();
            let mut app_dir = tauri::api::path::config_dir().unwrap();
            app_dir.push("read-chinese");
            println!("appdir is {}", app_dir.display());
            {
                #[cfg(debug_assertions)] // only include this code on debug builds
                window.open_devtools();
            }
            tauri::async_runtime::spawn(async move {
                let (mut rx, child) = Command::new_sidecar("app")
                    .expect("failed to setup `app` sidecar")
                    .args([app_dir.to_str().unwrap()]) // user_data.as_os_str())
                    .spawn()
                    .expect("Failed to spawn packaged node");

                let child = Arc::new(Mutex::new(child));
                window.listen_global("backend-in", move |event| {
                    println!("{}", event.payload().unwrap());
                    child
                        .lock()
                        .unwrap()
                        .write(format!("{}\n", event.payload().unwrap()).as_bytes())
                        .unwrap();
                });

                // read events such as stdout
                while let Some(event) = rx.recv().await {
                    if let CommandEvent::Stdout(line) = event {
                        window
                            .emit("message", Some(format!("'{}'", line)))
                            .expect("failed to emit event");
                        println!("{}", line)
                        // write to stdin
                        // child.write("message from Rust\n".as_bytes()).unwrap();
                    }
                }
            });

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
