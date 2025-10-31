import io
import grpc
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

            results = self.ocr.ocr(img_bytes, cls=True)

            text_list = []
            for line in results:
                for _, (text, confidence) in line:
                    text_list.append(text)

            return ocr_pb2.OCRResponse(text=text_list)

        except Exception as e:
            context.set_details(str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            return ocr_pb2.OCRResponse(text=[])
