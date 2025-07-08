FROM public.ecr.aws/lambda/python:3.11

RUN yum install -y clamav clamav-update && \
    freshclam && \
    mkdir -p /tmp/clamav && \
    chmod -R 777 /tmp/clamav

COPY app/ ${LAMBDA_TASK_ROOT}

CMD ["lambda_function.lambda_handler"]
