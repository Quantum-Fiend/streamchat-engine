use actix_web::{post, web, App, HttpServer, Responder, HttpResponse};
use serde::{Deserialize, Serialize};
use std::sync::Mutex;
use std::collections::HashMap;

// --- Models ---
#[derive(Deserialize, Debug)]
struct AnalyticsEvent {
    event_type: String, // "message_sent", "user_join"
    room_id: String,
    timestamp: i64,
}

struct AppState {
    // Simple in-memory counter for demo
    // RoomID -> MessageCount
    msg_counts: Mutex<HashMap<String, u64>>,
}

// --- Handlers ---
#[post("/ingest")]
async fn ingest_event(event: web::Json<AnalyticsEvent>, data: web::Data<AppState>) -> impl Responder {
    println!("Received Event: {:?}", event);
    
    // Process Event
    if event.event_type == "message" {
        let mut counts = data.msg_counts.lock().unwrap();
        let counter = counts.entry(event.room_id.clone()).or_insert(0);
        *counter += 1;
        println!("Room {} Message Count: {}", event.room_id, counter);
    }

    HttpResponse::Ok().json(serde_json::json!({"status": "recorded"}))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Starting Analytics Service (Rust) on port 9000...");
    
    let state = web::Data::new(AppState {
        msg_counts: Mutex::new(HashMap::new()),
    });

    HttpServer::new(move || {
        App::new()
            .app_data(state.clone())
            .service(ingest_event)
    })
    .bind(("127.0.0.1", 9000))?
    .run()
    .await
}
