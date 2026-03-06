<div dir="rtl">

# راهنمای کامل findns - اسکنر DNS Resolver

**فهرست مطالب:**

[1. findns چیست؟](#-findns-چیست-و-چه-کار-میکند) | [2. نصب](#-نصب-و-راهاندازی) | [🪟 ویندوز](#اجرا-روی-ویندوز-) | [3. دریافت لیست](#-دریافت-لیست-resolverها-fetch) | [4. اسکن کامل](#-اسکن-کامل-scan---دستور-اصلی) | [5. دستورات جداگانه](#-دستورات-جداگانه) | [6. Chain](#-پایپلاین-سفارشی-chain) | [7. فلگ‌ها](#%EF%B8%8F-فلگهای-عمومی) | [8. ورودی/خروجی](#-فرمت-ورودی-و-خروجی) | [9. سناریوها](#-سناریوهای-عملی) | [10. نکات](#-نکات-و-ترفندها)

---

## 1. findns چیست و چه کار می‌کند؟

findns یک ابزار خط فرمان است که DNS resolverها را تست می‌کند تا بفهمد کدام‌ها برای DNS tunneling (تانل DNS) مناسب هستند.

### DNS tunneling چیست؟

وقتی اینترنت فیلتر یا محدود است، می‌توان از پروتکل DNS برای عبور ترافیک استفاده کرد. ابزارهایی مثل DNSTT و Slipstream این کار را انجام می‌دهند. اما برای کار کردن، به یک DNS resolver نیاز دارید که:

- قابل دسترس باشد (ping جواب بدهد)
- واقعاً DNS resolve کند
- جواب جعلی (hijack) برنگرداند
- payload بزرگ DNS را پشتیبانی کند (EDNS)
- دامنه تانل شما را ببیند و resolve کند

findns همه این‌ها را به صورت خودکار تست می‌کند.

### چه پروتکل‌هایی پشتیبانی می‌شود؟

- **UDP DNS** (پورت 53) - روش کلاسیک
- **DoH** یعنی DNS-over-HTTPS (پورت 443) - شبیه ترافیک عادی HTTPS

### آیا به نصب dnstt یا slipstream نیاز دارم؟

**خیر!** findns به تنهایی تمام تست‌های DNS را انجام می‌دهد. فقط اگر بخواهید تست واقعی تانل (e2e) انجام دهید، به dnstt-client یا slipstream-client نیاز دارید. بدون آن‌ها هم اسکنر کامل کار می‌کند.

### dnstt-client چیست و چطور نصبش کنم؟

`dnstt-client` برنامه کلاینت پروژه [DNSTT](https://www.bamsoftware.com/software/dnstt/) است. این برنامه یک تانل DNS-over-UDP یا DNS-over-HTTPS درست می‌کند. findns از این برنامه برای **تست واقعی تانل** (e2e) استفاده می‌کند — یعنی واقعاً یک تانل می‌سازد و بررسی می‌کند اتصال برقرار می‌شود یا نه.

**نصب با Go (ساده‌ترین روش):**

</div>

```bash
go install www.bamsoftware.com/git/dnstt.git/dnstt-client@latest
```

<div dir="rtl">

**دانلود دستی:**

از [صفحه پروژه DNSTT](https://www.bamsoftware.com/software/dnstt/) باینری آماده دانلود کنید.

**بعد از دانلود، حتماً فایل را در PATH قرار دهید:**

</div>

```bash
# لینوکس/macOS:
sudo mv dnstt-client /usr/local/bin/
sudo chmod +x /usr/local/bin/dnstt-client

# یا PATH را به پوشه فعلی اضافه کنید:
export PATH=$PATH:$(pwd)
```

<div dir="rtl">

> **نکته مهم:** فقط گذاشتن فایل کنار findns روی لینوکس **کافی نیست** مگر اینکه پوشه فعلی در PATH باشد. روی ویندوز این مشکل وجود ندارد.

### کدام resolverها برای dnstt کار می‌کنند؟

بدون فلگ `--pubkey` هم findns بررسی می‌کند کدام resolverها **قابلیت** کار با تانل DNS را دارند:

- **resolve/tunnel**: بررسی می‌کند resolver می‌تواند NS record دامنه تانل شما را ببیند
- **edns**: بررسی می‌کند سایز payload بزرگ (1232 بایت) پشتیبانی می‌شود
- **nxdomain**: بررسی می‌کند resolver جواب جعلی نمی‌دهد

resolverهایی که همه این مراحل را پاس کنند، **با احتمال بالا** برای dnstt کار می‌کنند. فلگ `--pubkey` فقط تأیید نهایی (e2e) را اضافه می‌کند.

---

## 2. نصب و راه‌اندازی

### روش 1: دانلود باینری آماده (پیشنهادی)

</div>

**Linux (x64):**
```bash
curl -LO https://github.com/SamNet-dev/findns/releases/latest/download/findns-linux-amd64
chmod +x findns-linux-amd64
mv findns-linux-amd64 /usr/local/bin/findns
```

**Linux (ARM64):**
```bash
curl -LO https://github.com/SamNet-dev/findns/releases/latest/download/findns-linux-arm64
chmod +x findns-linux-arm64
mv findns-linux-arm64 /usr/local/bin/findns
```

**macOS (Intel):**
```bash
curl -LO https://github.com/SamNet-dev/findns/releases/latest/download/findns-darwin-amd64
chmod +x findns-darwin-amd64
mv findns-darwin-amd64 /usr/local/bin/findns
```

**macOS (Apple Silicon / M1/M2/M3):**
```bash
curl -LO https://github.com/SamNet-dev/findns/releases/latest/download/findns-darwin-arm64
chmod +x findns-darwin-arm64
mv findns-darwin-arm64 /usr/local/bin/findns
```

**Windows (x64):**

[findns-windows-amd64.exe](https://github.com/SamNet-dev/findns/releases/latest/download/findns-windows-amd64.exe)

<div dir="rtl">

بعد از نصب تست کنید:

</div>

```bash
findns --help
```

<div dir="rtl">

### روش 2: بیلد از سورس

نیازمند Go نسخه 1.24 یا بالاتر

</div>

```bash
git clone https://github.com/SamNet-dev/findns.git
cd findns
go build -o findns ./cmd
./findns --help
```

<div dir="rtl">

### روش 3: Go Install

</div>

```bash
go install github.com/SamNet-dev/findns/cmd@latest
```

<div dir="rtl">

### اجرا روی ویندوز 🪟

findns روی ویندوز **بدون نیاز به WSL یا لینوکس** کار می‌کند.

**دانلود مستقیم (بدون نصب چیزی):**

1. فایل [findns-windows-amd64.exe](https://github.com/SamNet-dev/findns/releases/latest/download/findns-windows-amd64.exe) را دانلود کنید
2. نام آن را به `findns.exe` تغییر دهید (اختیاری)
3. **cmd** یا **PowerShell** را در همان پوشه باز کنید (Shift + کلیک راست → Open PowerShell here)

**بیلد از سورس روی ویندوز:**

</div>

```powershell
git clone https://github.com/SamNet-dev/findns.git
cd findns
go build -o findns.exe ./cmd
```

<div dir="rtl">

**نحوه اجرا:** در تمام دستورات این راهنما به جای `findns` از `.\findns.exe` استفاده کنید:

</div>

```powershell
.\findns.exe fetch -o resolvers.txt
.\findns.exe scan -i resolvers.txt -o results.json --domain t.example.com
```

<div dir="rtl">

**نکات ویندوز:**
- **curl** از قبل در ویندوز 10/11 نصب است
- اگر ping فیل می‌شود → cmd را **Run as Administrator** باز کنید
- فایل‌های `dnstt-client.exe` و `slipstream-client.exe` را کنار `findns.exe` بگذارید
- در PowerShell برای ادامه دستورات طولانی از بک‌تیک `` ` `` استفاده کنید (به جای `\` در لینوکس)

---

## 3. دریافت لیست Resolverها (fetch)

قبل از اسکن، باید لیست resolver داشته باشید. دستور fetch به صورت خودکار از منابع عمومی دانلود می‌کند.

### دانلود resolverهای UDP جهانی

</div>

```bash
findns fetch -o resolvers.txt
```

<div dir="rtl">

این دستور از منبع trickest/resolvers حدود **17,000+** آی‌پی resolver عمومی را دانلود می‌کند.

### دانلود با resolverهای ایران (محلی)

</div>

```bash
findns fetch -o resolvers.txt --local
```

<div dir="rtl">

این دستور علاوه بر resolverهای جهانی، **7,800+** آی‌پی resolver ایرانی/منطقه‌ای از ir-resolvers را هم اضافه می‌کند. مجموعاً حدود **25,000+** resolver.

> **چرا resolverهای ایرانی مهم هستند؟** در شبکه ایران، resolverهای داخلی معمولاً سریع‌تر جواب می‌دهند و ممکن است محدودیت کمتری داشته باشند.

### دانلود resolverهای DoH

</div>

```bash
findns fetch -o doh-resolvers.txt --doh
```

<div dir="rtl">

این دستور آدرس‌های DoH (DNS-over-HTTPS) را جمع‌آوری می‌کند: 19 سرویس معروف (Google, Cloudflare, Quad9, AdGuard و ...) + لیست‌های عمومی از GitHub.

> فایل خروجی به صورت خودکار deduplicate می‌شود (تکراری‌ها حذف می‌شوند).

---

## 4. اسکن کامل (scan) - دستور اصلی

دستور scan مهم‌ترین و پیشنهادی‌ترین دستور است. تمام مراحل تست را به ترتیب اجرا می‌کند و فقط resolverهایی که همه مراحل را پاس کنند در خروجی نهایی می‌آیند.

### اسکن ساده (بدون دامنه تانل)

</div>

```bash
findns scan -i resolvers.txt -o results.json
```

<div dir="rtl">

مراحل: `ping -> resolve -> nxdomain`

این حالت بررسی می‌کند resolver زنده، فعال و بدون هایجک است. (برای رد کردن nxdomain از `--skip-nxdomain` استفاده کنید)

### اسکن کامل با دامنه تانل (پیشنهادی)

</div>

```bash
findns scan -i resolvers.txt -o results.json --domain t.example.com
```

<div dir="rtl">

مراحل: `ping -> resolve -> nxdomain -> edns -> resolve/tunnel`

### توضیح هر مرحله

**1. ping** — آیا سرور resolver از نظر شبکه قابل دسترس است؟ یک ICMP ping ارسال می‌کند و زمان پاسخ را اندازه می‌گیرد.
- متریک: `ping_ms` (میلی‌ثانیه)

**2. resolve** — آیا resolver واقعاً DNS resolve می‌کند؟ یک کوئری A record برای google.com ارسال می‌کند.
- متریک: `resolve_ms` (میلی‌ثانیه)

**3. nxdomain** — آیا resolver جواب جعلی می‌دهد (hijack)؟ یک دامنه تصادفی غیرموجود (مثل `nxd-abc123.invalid`) را کوئری می‌کند. resolver سالم باید NXDOMAIN برگرداند. resolver هایجک‌شده جواب NOERROR با آی‌پی جعلی برمی‌گرداند.
- متریک: `nxdomain_ok` (تعداد جواب‌های صحیح), `hijack` (0=سالم)

**4. edns** — resolver چه سایز payload DNS را پشتیبانی می‌کند؟ سایزهای 512, 900 و 1232 بایت تست می‌شود. هرچه بزرگ‌تر = تانل سریع‌تر.
- متریک: `edns_max` (بزرگ‌ترین سایز: 512, 900, یا 1232)

**5. resolve/tunnel** — آیا resolver دامنه تانل شما را می‌بیند؟ NS record و glue A record دامنه تانل را بررسی می‌کند. اگر resolver نتواند دامنه تانل را resolve کند، تانل کار نمی‌کند.
- متریک: `resolve_ms` (میلی‌ثانیه)

### اسکن با تست واقعی تانل DNSTT (اختیاری)

</div>

```bash
findns scan -i resolvers.txt -o results.json \
  --domain t.example.com --pubkey abc123def456...
```

<div dir="rtl">

مراحل: `ping -> resolve -> nxdomain -> edns -> resolve/tunnel -> e2e/dnstt`

نیازمند: `dnstt-client` و `curl` در PATH. این مرحله واقعاً dnstt-client را اجرا می‌کند، یک تانل SOCKS می‌سازد و با curl از طریق آن تانل یک صفحه وب را باز می‌کند.
- متریک: `e2e_ms` (کل زمان از شروع تا اتصال موفق)

### اسکن با تست واقعی Slipstream (اختیاری)

</div>

```bash
findns scan -i resolvers.txt -o results.json \
  --domain s.example.com --cert /path/to/cert.pem
```

<div dir="rtl">

نیازمند: `slipstream-client` و `curl` در PATH

### اسکن DoH

</div>

```bash
findns scan -i doh-resolvers.txt -o results.json --domain t.example.com --doh
```

<div dir="rtl">

مراحل: `doh/resolve -> doh/resolve/tunnel`

اسکن DoH با تست e2e:

</div>

```bash
findns scan -i doh-resolvers.txt -o results.json \
  --domain t.example.com --pubkey abc123... --doh
```

<div dir="rtl">

مراحل: `doh/resolve -> doh/resolve/tunnel -> doh/e2e`

### فلگ‌های دستور scan

| فلگ | توضیح | پیش‌فرض |
|-----|-------|---------|
| `--domain` | دامنه تانل (فعال‌سازی تست tunnel/edns) | — |
| `--pubkey` | کلید عمومی سرور DNSTT (فعال‌سازی تست e2e) | — |
| `--cert` | مسیر فایل گواهی Slipstream | — |
| `--test-url` | آدرسی که از طریق تانل تست شود | `https://httpbin.org/ip` |
| `--doh` | حالت DoH به جای UDP | `false` |
| `--skip-ping` | رد کردن مرحله ping (مفید اگر ICMP مسدود باشد) | `false` |
| `--skip-nxdomain` | رد کردن بررسی هایجک | `false` |
| `--top` | تعداد نتایج برتر در خروجی ترمینال | `10` |

---

## 5. دستورات جداگانه

هر مرحله از اسکن را می‌توانید به تنهایی هم اجرا کنید:

### ping - بررسی دسترسی‌پذیری

</div>

```bash
findns ping -i resolvers.txt -o ping-results.json
findns ping -i resolvers.txt -o ping-results.json -c 5 -t 2
```

<div dir="rtl">

`-c 5` = پنج بار ping بزن (پیش‌فرض: 3) | `-t 2` = تایم‌اوت 2 ثانیه (پیش‌فرض: 3)

### resolve - تست DNS Resolution

</div>

```bash
findns resolve -i resolvers.txt -o resolve-results.json --domain google.com
```

<div dir="rtl">

### resolve tunnel - بررسی NS Delegation

</div>

```bash
findns resolve tunnel -i resolvers.txt -o tunnel-results.json --domain t.example.com
```

<div dir="rtl">

بررسی می‌کند آیا resolver می‌تواند NS record دامنه تانل و glue A record آن را ببیند.

### nxdomain - تشخیص هایجک DNS

</div>

```bash
findns nxdomain -i resolvers.txt -o nxd-results.json
```

<div dir="rtl">

دامنه‌های تصادفی غیرموجود را کوئری می‌کند. resolver سالم: NXDOMAIN برمی‌گرداند. resolver هایجک‌شده: NOERROR با آی‌پی جعلی برمی‌گرداند.

### edns - تست سایز Payload

</div>

```bash
findns edns -i resolvers.txt -o edns-results.json --domain t.example.com
```

<div dir="rtl">

سایزهای 512, 900 و 1232 بایت را تست می‌کند.

### e2e dnstt - تست واقعی تانل DNSTT

</div>

```bash
findns e2e dnstt -i resolvers.txt -o e2e-results.json \
  --domain t.example.com --pubkey abc123...
```

<div dir="rtl">

نیازمند: `dnstt-client` و `curl`

### e2e slipstream - تست واقعی تانل Slipstream

</div>

```bash
findns e2e slipstream -i resolvers.txt -o e2e-results.json \
  --domain s.example.com --cert /path/to/cert.pem
```

<div dir="rtl">

نیازمند: `slipstream-client` و `curl`

### doh resolve - تست DoH Resolution

</div>

```bash
findns doh resolve -i doh-resolvers.txt -o doh-results.json --domain google.com
```

<div dir="rtl">

### doh resolve tunnel - تست DoH NS Delegation

</div>

```bash
findns doh resolve tunnel -i doh-resolvers.txt -o doh-tunnel-results.json \
  --domain t.example.com
```

<div dir="rtl">

### doh e2e - تست واقعی تانل از طریق DoH

</div>

```bash
findns doh e2e -i doh-resolvers.txt -o doh-e2e-results.json \
  --domain t.example.com --pubkey abc123...
```

<div dir="rtl">

نیازمند: `dnstt-client` و `curl`

---

## 6. پایپلاین سفارشی (chain)

دستور chain به شما اجازه می‌دهد مراحل دلخواه را ترکیب کنید. فقط resolverهایی که هر مرحله را پاس کنند به مرحله بعد می‌روند.

**مثال ساده:**

</div>

```bash
findns chain -i resolvers.txt -o results.json \
  --step "ping" \
  --step "resolve:domain=google.com"
```

<div dir="rtl">

**مثال کامل:**

</div>

```bash
findns chain -i resolvers.txt -o results.json \
  --step "ping:count=1" \
  --step "resolve:domain=google.com,count=1" \
  --step "nxdomain:count=2" \
  --step "edns:domain=t.example.com" \
  --step "resolve/tunnel:domain=t.example.com" \
  --step "e2e/dnstt:domain=t.example.com,pubkey=abc123,timeout=10"
```

<div dir="rtl">

فرمت هر step: `type:key=val,key=val`

**پارامترهای مشترک:**
- `count=N` — تعداد تلاش (پیش‌فرض: مقدار فلگ `-c`)
- `timeout=N` — تایم‌اوت به ثانیه (پیش‌فرض: مقدار فلگ `-t`)

### لیست تمام stepها

| Step | پارامترهای لازم | متریک خروجی |
|------|----------------|-------------|
| `ping` | — | `ping_ms` |
| `resolve` | `domain` | `resolve_ms` |
| `resolve/tunnel` | `domain` | `resolve_ms` |
| `nxdomain` | — | `hijack`, `nxdomain_ok` |
| `edns` | `domain` | `edns_max` |
| `e2e/dnstt` | `domain`, `pubkey` | `e2e_ms` |
| `e2e/slipstream` | `domain` | `e2e_ms` |
| `doh/resolve` | `domain` | `resolve_ms` |
| `doh/resolve/tunnel` | `domain` | `resolve_ms` |
| `doh/e2e` | `domain`, `pubkey` | `e2e_ms` |

**مثال DoH chain:**

</div>

```bash
findns chain -i doh-resolvers.txt -o results.json \
  --step "doh/resolve:domain=google.com" \
  --step "doh/resolve/tunnel:domain=t.example.com" \
  --step "doh/e2e:domain=t.example.com,pubkey=abc123"
```

<div dir="rtl">

---

## ⚙️ فلگ‌های عمومی

این فلگ‌ها روی همه دستورات کار می‌کنند:

| فلگ | مخفف | توضیح | پیش‌فرض |
|-----|------|-------|---------|
| `--input` | `-i` | فایل ورودی (متنی یا JSON) | الزامی |
| `--output` | `-o` | فایل خروجی JSON | الزامی |
| `--timeout` | `-t` | تایم‌اوت هر تلاش (ثانیه) | `3` |
| `--count` | `-c` | تعداد تلاش برای هر IP | `3` |
| `--workers` | — | تعداد workerهای موازی | `50` |
| `--e2e-timeout` | — | تایم‌اوت تست‌های e2e (ثانیه) | `10` |
| `--include-failed` | — | IPهای فیل‌شده از ورودی JSON را هم اسکن کن | `false` |

**تنظیم workers:**
- سرور ضعیف یا اینترنت کند: `--workers 20`
- سرور قوی: `--workers 100`
- پیش‌فرض 50 برای اکثر سرورها مناسب است

**تنظیم timeout:**
- شبکه ایران (resolverهای کند): `-t 5`
- سرور خارجی (پاسخ سریع): `-t 2`

---

## 8. فرمت ورودی و خروجی

### ورودی (Input)

**حالت 1: فایل متنی ساده** (یک آی‌پی یا URL در هر خط)

</div>

```
8.8.8.8
1.1.1.1
9.9.9.9
# این یک کامنت است (نادیده گرفته می‌شود)
```

<div dir="rtl">

برای DoH:

</div>

```
https://dns.google/dns-query
https://cloudflare-dns.com/dns-query
```

<div dir="rtl">

**حالت 2: خروجی JSON از اسکن قبلی**

خروجی هر اسکن می‌تواند ورودی اسکن بعدی باشد! به صورت پیش‌فرض فقط IPهای passed (موفق) استفاده می‌شوند. با فلگ `--include-failed` همه IPها دوباره تست می‌شوند.

</div>

```bash
findns ping -i resolvers.txt -o step1.json
findns resolve -i step1.json -o step2.json --domain google.com
```

<div dir="rtl">

### خروجی (Output)

فایل JSON با این ساختار:

</div>

```json
{
  "steps": [
    {
      "name": "ping",
      "tested": 25616,
      "passed": 20480,
      "failed": 5136,
      "duration_secs": 45.2
    }
  ],
  "passed": [
    {
      "ip": "1.1.1.1",
      "metrics": {
        "ping_ms": 4.2,
        "resolve_ms": 15.3,
        "edns_max": 1232,
        "nxdomain_ok": 3,
        "hijack": 0
      }
    }
  ],
  "failed": [
    {"ip": "9.9.9.9"}
  ]
}
```

<div dir="rtl">

- **steps:** خلاصه هر مرحله (چند تا تست شد، چند تا پاس شد)
- **passed:** لیست resolverهای موفق با متریک‌ها (مرتب شده بر اساس عملکرد)
- **failed:** لیست resolverهای ناموفق

---

## 9. سناریوهای عملی

### سناریو 1: پیدا کردن بهترین resolver UDP برای DNSTT

</div>

```bash
# مرحله 1 - دانلود resolverها (با لیست ایران)
findns fetch -o resolvers.txt --local

# مرحله 2 - اسکن کامل
findns scan -i resolvers.txt -o results.json --domain t.mysite.com

# مرحله 3 - استفاده در dnstt-client
dnstt-client -udp BEST_IP:53 -pubkey-file server.pub t.mysite.com 127.0.0.1:1080
```

<div dir="rtl">

نتایج به صورت TUI نمایش داده می‌شود و در results.json ذخیره می‌شود. اولین IP در لیست passed بهترین resolver است.

### سناریو 2: پیدا کردن resolver DoH برای DNSTT

</div>

```bash
# مرحله 1 - دانلود لیست DoH
findns fetch -o doh.txt --doh

# مرحله 2 - اسکن DoH
findns scan -i doh.txt -o doh-results.json --domain t.mysite.com --doh

# مرحله 3 - استفاده
dnstt-client -doh BEST_URL -pubkey-file server.pub t.mysite.com 127.0.0.1:1080
```

<div dir="rtl">

### سناریو 3: اسکن سریع (فقط ping + resolve)

</div>

```bash
findns scan -i resolvers.txt -o results.json --skip-nxdomain
```

<div dir="rtl">

### سناریو 4: اسکن وقتی ICMP مسدود است

</div>

```bash
findns scan -i resolvers.txt -o results.json \
  --domain t.mysite.com --skip-ping
```

<div dir="rtl">

### سناریو 5: فیلتر چندمرحله‌ای با chain

</div>

```bash
findns chain -i resolvers.txt -o results.json \
  --step "ping:count=1" \
  --step "resolve:domain=google.com,count=1" \
  --step "nxdomain:count=2" \
  --step "edns:domain=t.mysite.com" \
  --step "resolve/tunnel:domain=t.mysite.com"
```

<div dir="rtl">

مزیت: مرحله اول (`ping:count=1`) خیلی سریع فیلتر می‌کند و مراحل بعدی فقط روی resolverهای زنده اجرا می‌شوند.

### سناریو 6: استفاده از خروجی یک اسکن در اسکن بعدی

</div>

```bash
# مرحله 1 - فقط ping
findns ping -i resolvers.txt -o alive.json

# مرحله 2 - resolve فقط روی resolverهای زنده
findns resolve -i alive.json -o resolved.json --domain google.com

# مرحله 3 - nxdomain فقط روی resolverهای کارآمد
findns nxdomain -i resolved.json -o clean.json
```

<div dir="rtl">

هر مرحله فقط IPهای "passed" از مرحله قبل را تست می‌کند.

### سناریو 7: تست فقط resolverهای ایرانی

</div>

```bash
echo "10.202.10.10" > my-resolvers.txt
echo "10.202.10.11" >> my-resolvers.txt
echo "85.15.1.14" >> my-resolvers.txt

findns scan -i my-resolvers.txt -o results.json --domain t.mysite.com
```

<div dir="rtl">

### سناریو 8: تست با تعداد worker کمتر (سرور ضعیف)

</div>

```bash
findns scan -i resolvers.txt -o results.json \
  --domain t.mysite.com --workers 10 -t 5
```

<div dir="rtl">

---

## 10. نکات و ترفندها

**نکته 1: سرعت اسکن**
25,000 resolver با 50 worker حدود 5-15 دقیقه طول می‌کشد (بسته به شبکه). با `--workers 100` سریع‌تر می‌شود اما بار بیشتری روی سرور می‌گذارد.

**نکته 2: مرتب‌سازی نتایج**
نتایج بر اساس آخرین مرحله مرتب می‌شوند:
- اگر آخرین مرحله edns باشد: بر اساس `edns_max`
- اگر آخرین مرحله resolve/tunnel باشد: بر اساس `resolve_ms`
- اگر آخرین مرحله e2e باشد: بر اساس `e2e_ms`

**نکته 3: کجا اجرا کنیم؟**
بهترین جا یک سرور VPS است (نه کامپیوتر شخصی). چون سرور اینترنت پایدار و سریع دارد. می‌توانید روی همان سروری که DNSTT server دارید اجرا کنید.

**نکته 4: --top**
به صورت پیش‌فرض 10 نتیجه برتر در ترمینال نمایش داده می‌شود. برای دیدن بیشتر: `findns scan ... --top 50`. تمام نتایج همیشه در فایل JSON ذخیره می‌شوند.

**نکته 5: edns_max چقدر مهم است؟**
- `512`: حداقل (تانل کند)
- `900`: خوب
- `1232`: عالی (سریع‌ترین تانل)

resolverهایی با `edns_max=1232` بهترین انتخاب هستند.

**نکته 6: هایجک چیست و چرا مهم است؟**
بعضی ISPها و resolverها وقتی دامنه‌ای وجود ندارد، به جای NXDOMAIN شما را به صفحه تبلیغاتی یا صفحه خطای خودشان هدایت می‌کنند. این resolverها ممکن است تانل DNS را خراب کنند.

**نکته 7: فرق scan و chain**
- `scan`: خودکار مراحل را تنظیم می‌کند. ساده‌تر است.
- `chain`: شما مراحل را دستی تعریف می‌کنید. انعطاف بیشتر.

برای اکثر کاربران scan کافی است.

**نکته 8: اگر خطای "permission denied" گرفتید**
ping نیاز به دسترسی خاص دارد: `sudo findns scan ...` یا از `--skip-ping` استفاده کنید.

**نکته 9: اگر هیچ resolver پاس نشد**
- timeout را افزایش دهید: `-t 5` یا `-t 10`
- count را کم کنید: `-c 1`
- `--skip-nxdomain` امتحان کنید
- `--skip-ping` امتحان کنید
- مطمئن شوید دامنه تانل درست تنظیم شده (NS record)

**نکته 10: DoH یا UDP؟**

| | UDP (پورت 53) | DoH (پورت 443) |
|---|---|---|
| سرعت | سریع‌تر | کندتر |
| تعداد resolver | بیشتر | کمتر |
| قابل شناسایی | بله (DPI) | سخت (شبیه HTTPS) |
| مسدود شدن | ممکن | سخت‌تر |

> پیشنهاد: اول UDP امتحان کنید. اگر کار نکرد، DoH بزنید.

</div>
