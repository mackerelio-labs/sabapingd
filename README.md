# sabapingd

**sabapingd**は、ネットワーク機器などに対してPingを試行し、パケットロス率とRound-Trip Timeを計測して、Mackerelにメトリックとして送るプログラムです。

## 特徴

- 複数台に対しPingを試行し、指定したホストのメトリックとして投稿します
- ホストの指定には、ホストID または custom identifier を利用することができます
- Mackerelとの通信が途絶えた場合でもプログラム内部でキャッシュし、通信が再開したときに一斉に送信します

## 使い方

1. `config.yaml.sample を` `config.yaml` という名前でコピーします
1. `config.yaml` を開き加工します
1. `sabapingd -config config.yaml` で起動します

## 設定ファイルの内容

```yaml
x-api-key: xxxxx
# disk-cache: # 通信が長時間途絶えた場合にファイルに未送信データを書き出します
#   directory: cache
#   size: 10MB
# privileged: false # ICMP による Ping を利用する場合は、 true に指定してください
                    # 特権ユーザーでの起動もしくは、実行バイナリに capability を付与する等の措置が必要です
collector:
- host-id: xxxxx
  host: 192.2.0.1 # IPアドレス
  # custom-identifier: xxxx # host-id の代わりに利用できます
  # average: true # 平均メトリックのみを送信する場合、true に指定してください
```
