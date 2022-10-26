#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

mod database;
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

// Learn more about Tauri commands at
// https://tauri.app/v1/guides/features/command
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello {}, {}! You've been greeted from Rust!", name, name)
}

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
    // sqlx::migrate!().run(<&your_pool OR &mut your_connection>).await?;
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            greet,
            send_message,
            database::learning_target
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
