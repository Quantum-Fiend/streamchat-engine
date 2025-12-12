from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import time

app = FastAPI()

class TextPayload(BaseModel):
    text: str
    user_id: str

class ModerationResult(BaseModel):
    is_toxic: bool
    confidence: float
    original_text: str
    filtered_text: str

# Mock dictionary for "bad words" to simulate AI
BAD_WORDS = {"spam", "fail", "bad"}

@app.post("/moderate", response_model=ModerationResult)
async def moderate_text(payload: TextPayload):
    # Simulate processing time of a heavy model
    # time.sleep(0.1) 
    
    text_lower = payload.text.lower()
    is_toxic = any(word in text_lower for word in BAD_WORDS)
    
    # In a real scenario, we would allow the text but maybe flag it.
    # Or return star-out version.
    filtered_text = payload.text
    if is_toxic:
        filtered_text = "*" * len(payload.text)

    return ModerationResult(
        is_toxic=is_toxic,
        confidence=0.95 if is_toxic else 0.0,
        original_text=payload.text,
        filtered_text=filtered_text
    )

@app.get("/health")
def health_check():
    return {"status": "AI Model Loaded", "model": "Mock-Transformer"}
