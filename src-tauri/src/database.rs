use serde::Serialize;
use sqlx::sqlite::SqliteConnection;
use sqlx::Connection;
use sqlx::QueryBuilder;
use std::vec::Vec;

#[derive(sqlx::FromRow, Serialize)]
pub struct Target {
    word: String,
    occurance: i32,
}

#[tauri::command]
pub async fn learning_target(book_ids: Vec<i32>) -> Result<Vec<Target>, String> {
    let mut conn =
        SqliteConnection::connect("sqlite:/home/dlangevi/.config/read-chinese/db.sqlite3")
            .await
            .map_err(|e| e.to_string())?;
    let mut builder: QueryBuilder<sqlx::Sqlite> = QueryBuilder::new(
"SELECT word, sum(count) as occurance FROM frequency 
WHERE NOT EXISTS (
    SELECT word
    FROM words
    WHERE words.word==frequency.word
) ",
    );
    if book_ids.len() == 0 {
        builder.push("AND book in ( ");
        let mut separated = builder.separated(", ");
        for value_type in book_ids.iter() {
            separated.push_bind(value_type);
        }
        separated.push_unseparated(") ");
    }
    builder.push("GROUP BY word
            ORDER BY occurance DESC
            LIMIT 200");
    let query = builder.build_query_as::<Target>().fetch_all(&mut conn)
        .await
        .map_err(|e| e.to_string())?;
    Ok(query)
}
