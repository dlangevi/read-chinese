use serde::ser::{Serialize, SerializeStruct, Serializer};
use sqlx::sqlite::SqliteConnection;
use sqlx::Connection;
use std::vec::Vec;

#[derive(sqlx::FromRow)]
pub struct Target {
    word: String,
    occurance: i32,
}

impl Serialize for Target {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        let mut state = serializer.serialize_struct("Target", 2)?;
        state.serialize_field("word", &self.word)?;
        state.serialize_field("occurance", &self.occurance)?;
        state.end()
    }
}

#[tauri::command]
pub async fn learning_target(book_ids: Vec<i32>) -> Result<Vec<Target>, String> {
    let attempt =
        SqliteConnection::connect("sqlite:/home/dlangevi/.config/read-chinese/db.sqlite3").await;
    match attempt {
        Result::Ok(mut conn) => {
            let query_str = match book_ids.len() {
                0 => sqlx::query_as::<_, Target>(
                    "SELECT word, sum(count) as occurance FROM frequency 
                WHERE NOT EXISTS (
                    SELECT word
                    FROM words
                    WHERE words.word==frequency.word
                )
                GROUP BY word
                ORDER BY occurance DESC
                LIMIT 200",
                ),
                _ => {
                    let book_ids_string =
                        book_ids.iter().map(|id| id.to_string()).collect::<Vec<_>>();
                    let book_ids_joined = book_ids_string.join(",");

                    sqlx::query_as::<_, Target>(
                        "SELECT word, sum(count) as occurance FROM frequency 
                    WHERE NOT EXISTS (
                        SELECT word
                        FROM words
                        WHERE words.word==frequency.word
                    )
                    AND book in (?1)
                    GROUP BY word
                    ORDER BY occurance DESC
                    LIMIT 200",
                    )
                    .bind(book_ids_joined)
                }
            };
            let query = query_str.fetch_all(&mut conn).await;

            match query {
                Result::Ok(recs) => {
                    // Ok("success".to_string())
                    Ok(recs)
                }
                Result::Err(_err) => Err("Failed to query".to_string()),
            }
        }
        Result::Err(_err) => Err("Failed To Connect".to_string()),
    }
}
