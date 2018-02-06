Hatch is Another Try at Configuration Helper

To use it define a config struct which satisfies the `Hatchling` interface,
then call `hatch.NewWithConfig("application", "yaml", ".", "..", "...").Unmarshal(config)`

Check the test to see how it can be used and what are all the supported tag fields.

Tag fields syntax

```
hatch:"required,default"

e.g.

hatch:"true" -> required field
hatch:",hello world" -> optional field, with default "hello world"
hatch:",string1,string2,string3" -> optional field, with default "string1,string2,string3"

mapstructure:"env_name"
```
