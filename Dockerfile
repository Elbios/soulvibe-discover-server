FROM mcr.microsoft.com/dotnet/runtime-deps:8.0-alpine
WORKDIR /app

# the ELF produced above (name may differ if your project name does)
COPY --chmod=755 soulvibe_server /app/soulvibe_server
COPY --chmod=755 spotseek /app/spotseek
COPY --chmod=777 appsettings.json /app/appsettings.json

# static assets for your Go server
COPY templates/ /app/templates/
COPY static/ /app/static/

ENV DOTNET_gcServer=0 \
    COMPlus_GCHeapHardLimit=268435456

EXPOSE 8080
USER 101
CMD ["/app/soulvibe_server"]
