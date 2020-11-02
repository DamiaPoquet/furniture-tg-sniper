FROM chromedp/headless-shell:latest

RUN apt-get update && apt-get install -y ca-certificates

COPY ./furniture-tg-sniper /usr/local/bin/furniture-tg-sniper

ENTRYPOINT ["/usr/local/bin/furniture-tg-sniper"]
