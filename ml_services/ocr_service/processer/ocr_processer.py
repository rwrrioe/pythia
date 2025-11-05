import io
import grpc
from PIL import Image
import numpy as np
from paddleocr import PaddleOCR
from shared.gen.python.ocr import ocr_pb2, ocr_pb2_grpc


class OCRServiceServicer(ocr_pb2_grpc.OCRServiceServicer):
    def __init__(self):
        self.ocr = PaddleOCR(use_angle_cls=True, lang='de')

    def Recognize(self, request, context):
        try:
            with io.BytesIO(request.image_data) as img_buf:
                img = Image.open(img_buf).convert("RGB")
                img = np.array(img)

            results = self.ocr.predict(img)
    
            if not results:
                return ocr_pb2.OCRResponse(text=[])

            texts = results[0].get("rec_texts", [])


            return ocr_pb2.OCRResponse(text=texts)

        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return ocr_pb2.OCRResponse(text=[])
