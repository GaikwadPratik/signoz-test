FROM ubuntu:22.04

ARG BUILD_DATA=build-data
ENV APP_CONFIG=app-config.yml
ENV SVC_NAME=signoz-test

EXPOSE 5555

# Copy the binary file
COPY ${BUILD_DATA}/${SVC_NAME} /usr/local/bin/${SVC_NAME}
# Copy app config file
COPY ${APP_CONFIG} .


ENTRYPOINT [ "/bin/bash" ]
CMD [ "-c", "${SVC_NAME} server --config ./${APP_CONFIG}" ]