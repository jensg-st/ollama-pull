FROM gerke74/ollama-model-loader as downloader

RUN /ollama-pull gemma:2b

FROM ollama/ollama 

ENV OLLAMA_HOST "0.0.0.0"

COPY --from=downloader /root/.ollama /root/.ollama
