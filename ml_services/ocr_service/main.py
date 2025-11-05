import grpc
from concurrent import futures
from processer.ocr_processer import OCRServiceServicer
from shared.gen.python.ocr import ocr_pb2_grpc


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    ocr_pb2_grpc.add_OCRServiceServicer_to_server(OCRServiceServicer(), server)

    port = 50051
    server.add_insecure_port(f"[::]:{port}")

    print(f"OCR gRPC server started on port {port}")
    server.start()

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        print("\nShutting down OCR gRPC server gracefully...")
        server.stop(0)


if __name__ == "__main__":
    serve()
