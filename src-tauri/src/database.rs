use serde::Serialize;
use std::vec::Vec;
use tauri::State;

#[derive(sqlx::FromRow, Serialize)]
pub struct Target {
    word: String,
    occurance: i32,
}

#[tauri::command]
pub async fn learning_target(
    pool: State<'_, sqlx::SqlitePool>,
) -> Result<Vec<Target>, String> {
    let mut conn = pool.acquire().await.map_err(|e| e.to_string())?;
    let query = sqlx::query_as::<_, Target>(
        r#"
SELECT word, sum(count) as occurance FROM frequency 
WHERE NOT EXISTS (
    SELECT word
    FROM words
    WHERE words.word==frequency.word
) 
GROUP BY word
ORDER BY occurance DESC
LIMIT 200
        "#,
    )
    .fetch_all(&mut conn)
    .await
    .map_err(|e| e.to_string())?;
    Ok(query)
}

#[tauri::command]
pub async fn learning_target_book(
    pool: State<'_, sqlx::SqlitePool>,
    book_id: i32,
) -> Result<Vec<Target>, String> {
    let mut conn = pool.acquire().await.map_err(|e| e.to_string())?;
    let query = sqlx::query_as::<_, Target>(
        r#"
SELECT word, sum(count) as occurance FROM frequency 
WHERE NOT EXISTS (
    SELECT word
    FROM words
    WHERE words.word==frequency.word
) 
AND book == (?1)
GROUP BY word
ORDER BY occurance DESC
LIMIT 200
        "#,
    )
    .bind(book_id)
    .fetch_all(&mut conn)
    .await
    .map_err(|e| e.to_string())?;
    Ok(query)
}
