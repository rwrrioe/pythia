FROM nvidia/cuda:12.9.0-runtime-ubuntu22.04

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        python3 python3-pip libgl1 libglib2.0-0 git curl && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN pip3 install --upgrade pip && \
    pip install paddlepaddle-gpu==3.2.1 \
    -i https://www.paddlepaddle.org.cn/packages/stable/cu129/

COPY requirements.txt .

RUN pip3 install -r requirements.txt

COPY . .
ENV PYTHONPATH="${PYTHONPATH}:/app:/app/shared"
ENV PYTHONUNBUFFERED=1

CMD ["python3", "main.py"]

