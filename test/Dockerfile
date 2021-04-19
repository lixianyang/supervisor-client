FROM python:3.9.1

RUN pip install supervisor

COPY supervisord.conf /etc/supervisord.conf

ENTRYPOINT ["supervisord", "-c", "/etc/supervisord.conf"]
