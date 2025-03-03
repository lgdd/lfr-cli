FROM python:3.10-alpine

WORKDIR /usr/src/app

COPY pyproject.toml .
RUN pip install --no-cache-dir -e .

COPY . .

RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.domains
RUN sed -i -e 's/localhost/liferay/g' dxp-metadata/com.liferay.lxc.dxp.mainDomain

# add path to Liferay client extension metadata
ENV LIFERAY_ROUTES_DXP=/usr/src/app/dxp-metadata

CMD ["python", "./main.py"]