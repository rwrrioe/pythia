import io
import grpc
import json
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
                img_bytes = img_buf.read()

            results = self.ocr.predict(img_bytes, cls=True)
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
