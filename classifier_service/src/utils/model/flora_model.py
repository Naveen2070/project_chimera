import torch
from torchvision import transforms
import os


# Define class names
class_names = [
    "daisy",
    "dandelion",
    "lotus",
    "roses",
    "sunflowers",
    "tulips",
]  # ⬅️ Update if needed


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
