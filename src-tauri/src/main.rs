#![cfg_attr(
    all(not(debug_assertions), target_os = "windows"),
    windows_subsystem = "windows"
)]

use std::sync::{Arc, Mutex};
use std::collections::HashMap;

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
    // Lol there has got to be a better way of writing this
    let mut map = HashMap::new();
    map.insert("function", function);
    map.insert("args", args);
    let client = reqwest::Client::new();
    let resp = client.post("http://localhost:3451/ipc")
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
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![greet, send_message])
        .setup(|app| {
            let window = app.get_window("main").unwrap();
            // #[cfg(debug_assertions)] // only include this code on debug builds
            let mut app_dir = tauri::api::path::config_dir().unwrap();
            app_dir.push("read-chinese");
            println!("appdir is {}", app_dir.display());
            window.open_devtools();
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
