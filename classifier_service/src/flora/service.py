import torch

from src.utils.image.tools import decode_image
from src.utils.model.flora_model import load_latest_model
from src.utils.model.flora_model import class_names
from src.utils.model.flora_model import transform


class FloraService:
    def __init__(self, model_dir: str = "saved_models/flora/scripted/"):
        self.model = load_latest_model(model_dir)
        self.model.eval()

    def predict(self, image_base64: str) -> str:
        image = decode_image(image_base64)
        input_tensor = transform(image).unsqueeze(0)

        with torch.no_grad():
            outputs = self.model(input_tensor)
            _, predicted = torch.max(outputs, 1)
            predicted_label = class_names[predicted.item()]

        return predicted_label
