Monitooringusüsteemi OpenTelemetry rakendamise POC

Repos on eraldikäivitatavad rakendused:

`app1` - rakendus valmistab lihtsa OpenTelemetry mõõteriista
(instrument), teeb sellega sarja mõõtmisi ja väljastab mõõtmistulemused
standardväljundisse.

`app2` - on teegi `opentelemetry-go` näidisrakendus `basic`, kommenteeritud ja kohandatud.

Rakendused on vormistatud ühisesse Go-moodulisse. Mooduliteed:

- `github.com/e-gov/opentelemetry-poc/app1`
- `github.com/e-gov/opentelemetry-poc/app2`.

Eeldatakse, et POC-i rakendus ehitatakse ja seda käitatakse samas masinas. Rakendustes kasutatakse OpenTelemetry Go-teegi (`opentelemetry-go`) viimast GitHub-is avaldatud seisu. Kuna see seis ei tarvitse veel olla `go get` mehhanismiga kättesaadav, siis on moodulifailis (`go.mod`) seadistatud pöördumised OpenTelemetry Go-teegi poole eeldusega, et OpenTelemetry Go-teek (`opentelemetry-go`) on kloonitud POC-i ehitamise ja käitamise masinasse kausta `OpenTelemetry-POC` kõrvalkausta. (Vt `replace`-direktiivid `go.mod` failis).

Rakenduse käivitamiseks:

- taga, et OpenTelemetry Go-teek (`opentelemetry-go`) on kloonitud POC-i ehitamise ja käitamise masinasse kausta `OpenTelemetry-POC` kõrvalkausta.
- liigu rakenduse kausta
- `go run .`
