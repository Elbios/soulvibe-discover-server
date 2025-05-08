FROM mcr.microsoft.com/dotnet/runtime-deps:8.0-alpine
WORKDIR /app

# the ELF produced above (name may differ if your project name does)
COPY soulvibe_server /app/soulvibe_server
COPY spotseek /app/spotseek

# static assets for your Go server
COPY templates/ /app/templates/
COPY static/ /app/static/

ENV DOTNET_gcServer=0 
#\
    #COMPlus_GCHeapHardLimit=134217728   # 128 MiB hard cap, tune as you like

EXPOSE 8080
USER 101
CMD ["/app/soulvibe_server"]