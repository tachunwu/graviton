nats request '$GVTN.KV.key1' -H "Ops:set" '{"Value":"Value 1"}' -s localhost:4322
nats request '$GVTN.KV.key2' -H "Ops:set" '{"Value":"Value 2"}' -s localhost:4322
nats request '$GVTN.KV.key3' -H "Ops:set" '{"Value":"Value 3"}' -s localhost:4322

nats request '$GVTN.KV.key1' '' -H "Ops:get"  -s localhost:4322
nats request '$GVTN.KV' '' -H "Ops:scan" -H "Start-Key:key1" -H "End-Key:key3" -s localhost:4322

nats request '$GVTN.KV.key2' '' -H "Ops:delete" -s localhost:4322
nats request '$GVTN.KV.key2' '' -H "Ops:get" -s localhost:4322


nats request '$GVTN.KV.key{{count}}' -H "Ops:set" '{"Value":"Value {{count}}"}' -s localhost:4322 --count 10000
