FROM nvidia/cuda:12.2.0-runtime-ubuntu22.04

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        python3 python3-pip libgl1 libglib2.0-0 git curl && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .

RUN pip3 install --upgrade pip && \
    pip3 install --no-cache-dir -r requirements.txt && \
    pip3 install --no-cache-dir paddlepaddle-gpu==2.5.2.post120 -i https://mirror.baidu.com/pypi/simple && \
    pip3 install --no-cache-dir paddleocr==2.7.0.3


COPY . .

ENV PYTHONUNBUFFERED=1

CMD ["python3", "main.py"]
