from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import base64
from io import BytesIO
from PIL import Image
import torch
from torchvision import transforms
import os

app = FastAPI()

# Define class names
class_names = [
    "daisy",
    "dandelion",
    "lotus",
    "roses",
    "sunflowers",
    "tulips",
]  # ⬅️ Update if needed


# Define input model
class ImagePayload(BaseModel):
    image_base64: str


# Define transform
transform = transforms.Compose(
    [
        transforms.Resize((224, 224)),
        transforms.ToTensor(),
        transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
    ]
)


# Load latest scripted model
def load_latest_model(model_dir="models"):
    files = [
        f
        for f in os.listdir(model_dir)
        if f.startswith("flora_model_script_v") and f.endswith(".pt")
    ]
    versions = [int(f.split("_v")[-1].split(".pt")[0]) for f in files]
    if not versions:
        raise RuntimeError("No model found.")

    latest_version = max(versions)
    model_path = os.path.join(model_dir, f"flora_model_script_v{latest_version}.pt")

    model = torch.jit.load(model_path, map_location="cpu")
    model.eval()
    return model


model = load_latest_model(model_dir="saved_models/flora/scripted/")


# Convert Base64 to PIL image
def decode_image(base64_str: str) -> Image.Image:
    try:
        img_data = base64.b64decode(base64_str)
        return Image.open(BytesIO(img_data)).convert("RGB")
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid image format")


@app.post("/predict")
def predict(payload: ImagePayload):
    image = decode_image(payload.image_base64)
    input_tensor = transform(image).unsqueeze(0)  # Add batch dim

    with torch.no_grad():
        outputs = model(input_tensor)
        _, predicted = torch.max(outputs, 1)
        predicted_label = class_names[predicted.item()]

    return {"prediction": predicted_label}
