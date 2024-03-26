# ollama-pull

This is an example repository to show how to build an Ollama model downloader for Docker builds. This tool can download Ollama models during a Docker build. Although almost all functions in the Ollama are exported the dependecies are hard to get right. 

The main logic of this little tool is in the server package. This source code is 99% Ollama code which is MIT-licensed: https://github.com/ollama/ollama/blob/main/LICENSE


```
FROM gerke74/ollama-model-loader as downloader

RUN /ollama-pull gemma:2b

FROM ollama/ollama 

ENV OLLAMA_HOST "0.0.0.0"

COPY --from=downloader /root/.ollama /root/.ollama

```

