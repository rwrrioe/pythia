import io
import grpc
import json
from PIL import Image
import numpy as np
from paddleocr import PaddleOCR
from shared.gen.python.ocr import ocr_pb2, ocr_pb2_grpc


class OCRServiceServicer(ocr_pb2_grpc.OCRServiceServicer):
    def __init__(self):
        self.ocr = PaddleOCR(lang='en')

    def Recognize(self, request, context):
        try:
            image_data = request.image_data
            lang = request.lang or 'en'

            if lang != self.ocr.lang:
                self.ocr = PaddleOCR(lang=lang)

            with io.BytesIO(image_data) as img_buf:
                img = Image.open(img_buf)
                img = np.array(img)

            results = self.ocr.predict(img, cls=True)
            if isinstance(results, list):
                data = results[0]
            else:
                data = results

            texts = data.get("rec_texts", [])
            return ocr_pb2.OCRResponse(text=texts)

        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return ocr_pb2.OCRResponse(text=[])
