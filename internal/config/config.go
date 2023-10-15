package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

type Config struct {
	LogLevel          string           `json:"logLevel,omitempty"`
	Ip                string           `json:"ip,omitempty"`
	Port              uint16           `json:"port,omitempty"`
	DebugPacketLogUid uint32           `json:"debugPacketLogUid,omitempty"`
	HttpPort          uint16           `json:"httpPort,omitempty"`
	TerrainCollect    bool             `json:"terrainCollect"`
	LuaShellFile      []string         `json:"luaShellFile"`
	Endpoints         *ConfigEndpoints `json:"endpoints,omitempty"`
	Protocols         *ConfigProtocols `json:"protocols,omitempty"`
	Keys              *ConfigKeys      `json:"keys,omitempty"`
}

type ConfigConsole struct {
	Enabled      bool   `json:"enabled,omitempty"`
	MuipEndpoint string `json:"muipEndpoint,omitempty"`
	MuipRegion   string `json:"muipRegion,omitempty"`
	MuipSign     string `json:"muipSign,omitempty"`
}

type ConfigEndpoints struct {
	MainEndpoint string              `json:"mainEndpoint,omitempty"`
	MainProtocol Protocol            `json:"mainProtocol,omitempty"`
	Console      *ConfigConsole      `json:"console,omitempty"`
	Mapping      map[Protocol]string `json:"mapping,omitempty"`
}

type ConfigProtocols struct {
	BaseProtocol Protocol            `json:"baseProtocol,omitempty"`
	Mapping      map[Protocol]string `json:"mapping,omitempty"`
}

type ConfigKeys struct {
	SharedKey  string            `json:"sharedKey,omitempty"`
	ServerKey  string            `json:"serverKey,omitempty"`
	ClientKeys map[uint32]string `json:"clientKeys,omitempty"`
}

var CONF *Config = nil

func GetConfig() *Config {
	return CONF
}

var FileNotExist = errors.New("config file not found")

func LoadConfig() error {
	filePath := "./config.json"
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}
	f, err := os.Open(filePath)
	if err != nil {
		return FileNotExist
	}
	defer func() {
		_ = f.Close()
	}()
	c := new(Config)
	d := json.NewDecoder(f)
	if err := d.Decode(c); err != nil {
		return err
	}
	if c.Endpoints == nil {
		return errors.New("no endpoint configured")
	}
	if c.Endpoints.Console == nil {
		c.Endpoints.Console = &ConfigConsole{}
	}
	if c.Protocols == nil {
		return errors.New("no protocol configured")
	}
	if c.Keys == nil {
		return errors.New("no key configured")
	}
	CONF = c
	return nil
}

var DefaultConfig = &Config{
	LogLevel:          "info",
	Ip:                "0.0.0.0",
	Port:              20045,
	DebugPacketLogUid: 100000001,
	HttpPort:          8080,
	Endpoints: &ConfigEndpoints{
		MainEndpoint: "{{ UPSTREAM_SERVER_ADDRESS }}",
		MainProtocol: "{{ UPSTREAM_SERVER_VERSION }}",
		Console: &ConfigConsole{
			Enabled:      false,
			MuipEndpoint: "http://{{ MUIP_SERVER_ADDRESS }}/api",
			MuipRegion:   "DEV_TianliPS",
			MuipSign:     "9H2UrJ5J4yZJf95FqMkqi628snEmzvyV9oAp",
		},
		Mapping: map[Protocol]string{
			"{{ CLIENT_VERSION }}": "{{ SERVICE_LISTEN_ADDRESS }}",
		},
	},
	Protocols: &ConfigProtocols{
		BaseProtocol: ProtocolMajor3Minor2,
		Mapping: map[Protocol]string{
			ProtocolMajor3Minor2: path.Join("data/mapping/", string(ProtocolMajor3Minor2)),
		},
	},
	Keys: defaultConfigKeys,
}

var defaultConfigKeys = &ConfigKeys{
	SharedKey: "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",
	// ServerKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAxbbx2m1feHyrQ7jP+8mtDF/pyYLrJWKWAdEv3wZrOtjOZzeL\nGPzsmkcgncgoRhX4dT+1itSMR9j9m0/OwsH2UoF6U32LxCOQWQD1AMgIZjAkJeJv\nFTrtn8fMQ1701CkbaLTVIjRMlTw8kNXvNA/A9UatoiDmi4TFG6mrxTKZpIcTInvP\nEpkK2A7Qsp1E4skFK8jmysy7uRhMaYHtPTsBvxP0zn3lhKB3W+HTqpneewXWHjCD\nfL7Nbby91jbz5EKPZXWLuhXIvR1Cu4tiruorwXJxmXaP1HQZonytECNU/UOzP6GN\nLdq0eFDE4b04Wjp396551G99YiFP2nqHVJ5OMQIDAQABAoIBAQDEeYZhjyq+avUu\neSuFhOaIU4/ZhlXycsOqzpwJvzEz61tBSvrZPA5LSb9pzAvpic+7hDH94jX89+8d\nNfO7qlADsVNEQJBxuv2o1MCjpCRkmBZz506IBGU60Kt1j5kwdCEergTW1q375z4w\nl8f7LmSL2U6WvKcdojTVxohBkIUJ7shtmmukDi2YnMfe6T/2JuXDDL8rvIcnfr5E\nMCgPQs+xLeLEGrIJdpUy1iIYZYrzvrpJwf9EJL3D0e7jkpbvAQZ8EF9YhEizJhOm\ndzTqW4PgW2yUaHYd3q5QjiILy7AC+oOYoTZln3RfjPOxl+bYjeMOWlqkgtpPQkAE\n4I64w8RZAoGBAPLR44pEkmTdfIIF8ZtzBiVfDZ29bT96J0CWXGVzp8x6bSu5J5jl\ns7sP8DEcjGZ6vHsLGOvkcNxzcnR3l/5HOz6TIuvVuUm36b1jHltq1xZStjGeKZs1\nihhJSu2lIA+TrK8FCRnKARJ0ughXGNZFItgeM230Sgjp2RL4ISXJ724XAoGBANBy\nS2RwNpUYvkCSZHSFnQM/jq1jldxw+0p4jAGpWLilEaA/8xWUnZrnCrPFF/t9llpb\ndTR/dCI8ntIMAy2dH4IUHyYKUahyHSzCAUNKpS0s433kn5hy9tGvn7jyuOJ4dk9F\no1PIZM7qfzmkdCBbX3NF2TGpzOvbYGJHHC3ssVr3AoGBANHJDopN9iDYzpJTaktA\nVEYDWnM2zmUyNylw/sDT7FwYRaup2xEZG2/5NC5qGM8NKTww+UYMZom/4FnJXXLd\nvcyxOFGCpAORtoreUMLwioWJzkkN+apT1kxnPioVKJ7smhvYAOXcBZMZcAR2o0m0\nD4eiiBJuJWyQBPCDmbfZQFffAoGBAKpcr4ewOrwS0/O8cgPV7CTqfjbyDFp1sLwF\n2A/Hk66dotFBUvBRXZpruJCCxn4R/59r3lgAzy7oMrnjfXl7UHQk8+xIRMMSOQwK\np7OSv3szk96hy1pyo41vJ3CmWDsoTzGs7bcdMl72wvKemRaU92ckMEZpzAT8cEMC\ncWKLb8yzAoGAMibG8IyHSo7CJz82+7UHm98jNOlg6s73CEjp0W/+FL45Ka7MF/lp\nxtR3eSmxltvwvjQoti3V4Qboqtc2IPCt+EtapTM7Wo41wlLCWCNx4u25pZPH/c8g\n1yQ+OvH+xOYG+SeO98Phw/8d3IRfR83aqisQHv5upo2Rozzo0Kh3OsE=\n-----END RSA PRIVATE KEY-----",
	ServerKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAvTMxK4jdNZ+3hBBYG7jFKvaPYE6vbpZOE9z8HkZztlftpp8v\nYLNE49Pr4mRKHSK40D4ozx0U+XobbYqU8xQZtAKi4KslSBWlxMZrY5LRjHSi2WeX\nornpS25RKkSjFvKf2rUIoFf72Tr+4VJK96QynXT6WRWM8pJbNgoXaQbA12IgQoB4\nb6epsAqYwy4lB1RMe4wP6NQbk4X7vJa0GTwofntAizDcJM5ULWjIFj+0pghd3PKp\nRftYGMwxOidtGvSnOPiBWRHdHURcZelxxcsVs5GLjhGvLZID0En6IkzUwDSghgTX\n+leiAy4Dt9dTu9j5VoZBXablufD6B0dbVdWC0QIDAQABAoIBAADRnQW81cU+w9Tt\n263CCNNo5gGzEXoLazrVht9MK4HHY2NIVtSGrhaE0yVjjPkPjohzLmrIC9BRFZeN\npRugs4RGiyUpBHQpiNI/UBVqHB9NHWvOsZauEosFpxMFYUXPBr9T732/J7s+0L3R\npPqpoqDxEcjLKbUeikfDdyq4kWLVO1grnt29Knpbm83j5a3RKLy9kEDpLbqe/C50\n983Lri4yBmfZSnmfTueD3GcnISQc6UYxp95N75j2lULt8mPCWbDjB39SQfYD8qcP\nm0xrLC2E5jhZktH5DFOXeFrNv5d1cqjWBCNEJ6Dw4SujO1cxtoTMXV1ToKdnVJn5\nF02WYsECgYEA8I5vwun+6FsKuOq9AeFYM2B6lcsalH6z0r8f6ypbOqTKH+ifdZRk\nQAOwI5JOCzBa/bIkdVl3zGcAcS7X9jDJh7lhHdss8fIcb2UlaVZoCHhEtJXQ5gTH\n32FQLKtyN1QkiUps9rzPd2F4SEJNWwVw1fPoyZJpbYszYU10+Rso76ECgYEAyVi1\nZo7+pRq0IqOSIfCJ+6LdIa83pxQtKZtIvhbyGaAw+tb3ja4Oqpr72PdylNHO361o\nQ/sh0JvbkYyIYPQbJZo+sysTfHdlJqYLzb5HGMrbWPHlPnSdMYIASDJGrtX4Smy9\nP6VK4nM38dNFnoTw6mX+AxKqXPauP27YVwLhhTECgYEA6pLK0vPg+W2F+CoXIxU+\nL+NdxmImyjT/X3u2QVitW3NEEneBv2NzmqS+BwHtDqYZpJgpSzFyS6UJXlVCjLSo\nYKxZ0oZevpPMPKgSIjT6/39f6ATLjvGMgfxf9R8+Ikvv0Nz9gmE9ofkvFK9qxV55\n2HifQKiAHC0IblLcxOlCMuECgYB7koQSo6RJdHAl6jnftp8Y30XUTJNdaZamOHWW\npMKFU7l72b8pJzA9KM10xbl++J18zhJ11oVUYLOLSrLQvkCC/X2JvOBCvYxJAhOw\nfB1qa+XfWuaVREDNh7nglWqoFw5BrycfDrU88fXd5wqNVY3+bgZNoIEKeSNMLx17\nmXsLoQKBgBk6E9b7i5fk2acud3EoUe6cNRw8ItMqGwxsoR33jkAdVSMXGbOH5tjX\nPvM8ZHmMI0xvPi3qZhlZ8cx2lXaS9hZovwlL9s0UGzYHml2gsiMDqunCp3770T7H\ntT3fryVzT5L3e1522rJx405tV/j6IkxifH2bbPfO3I/IbPwoeEYx\n-----END RSA PRIVATE KEY-----",

	ClientKeys: map[uint32]string{
		2: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAz/fyfozlDIDWG9e3Lb29+7j3c66wvUJBaBWP10rB9HTE6prj\nfcGMqC9imr6zAdD9q+Gr1j7egvqgi3Da+VBAMFH92/5wD5PsD7dX8Z2f4o65Vk2n\nVOY8Dl75Z/uRhg0Euwnfrved69z9LG6utmlyv6YUPAflXh/JFw7Dq6c4EGeR+Kej\nFTwmVhEdzPGHjXhFmsVt9HdXRYSf4NxHPzOwj8tiSaOQA0jC4E4mM7rvGSH5GX6h\nma+7pJnl/5+rEVM0mSQvm0m1XefmuFy040bEZ/6O7ZenOGBsvvwuG3TT4FNDNzW8\nDw9ExH1l6NoRGaVkDdtrl/nFu5+a09Pm/E0ElwIDAQABAoIBAQCtH17Cck+KJQYX\nj29xqG4qykNUDawbILiKCMkBE7553Wq/UcjmuuR4bVnML8ucS3mgR/BgHV3l8vUK\nnxvqRx/oGZkWNazbiuwL+ThAblLWqrEmYuZVCoQcAnvkT8tIqDWz7fhDEuZnnkMz\nZcATIZzgZUSa5IfP3u3rP+MrVbyaCdzJEeI0Yrv1XT+M5ddkKQrYgqC5kRiYi/Lj\nNcLJhqSVt8p37CdJx1PGHFjKKb4MZpANlNRgeTtWpGVfS0PJLzaiI1NyPSJv7xWZ\ngVhbK9+wQxqSG6KmZ4vpEvRI1zKiov5BsAFN+GfuD5mpn1Xo9CpzTfj/sO13VpHH\n+Mt80+yBAoGBAPYXVEcXug5zqkqXup4dp1S05saz1zWPhUhQm+CrbhgeTqpjngJJ\nEB79qMrGmyki0P/cGtbTcrHf8+i7gDlIGW0OMb4/jn4f5ACVD00iyvkHSGPn0Aim\nMoNOMbkGot7SkSnncwxXdawwDyTu2dofXuBr72+GYqgRAG52IuA0C0pRAoGBANhX\np/UyW/htB27frKch/rTKQKm12kBV20AkkRUQUibiiQyWueWKs+5bVaW5R5oDIhWx\nqftJtnEFWUvWaTHpHsB/bpjS3CJ6WknqNbpa3QIScpV1uw8V+Etz/K2/ftjyZzFo\nnqc+Jud5364xFdIlOsRj9gZnK83Wcui6EFxAer5nAoGBAJzTzzSjLUHqejqhKR98\nnFeCFZPJpjuO5AxqunvaJAYgwlcZtueT8j8dvgTDvrvfYTu85CnFhNFQfFrzqspW\nZUW3hwHL9R3xatboJ2Er7Bf5iSuJ3my0pXpCSbO1Q/QmUrZWtl3GGsqJsg0CXjkA\nRvFUN7ll9ddPRmwewykIYa2RAoGAcmKuWFNfE1O6WWIELH4p6LcDR3fyRI/gk+KB\nnyx48zxVkAVllrsmdYFvIGd9Ny4u6F9+a3HG960HULS1/ACxFMCL3lumrsgYUvp1\nm+mM7xqH4QRVeh14oZRa5hbY36YS76nMMMsI0Ny8aqJjUjADCXF81FfabkPTj79J\nBS3GeEMCgYAXmFIT079QokHjJrWz/UaoEUbrNkXB/8vKiA4ZGSzMtl3jUPQdXrVf\ne0ofeKiqCQs4f4S0dYEjCv7/OIijV5L24mj/Z1Q4q++S5OksKLPPAd3gX4AYbRcg\nPS4rUKl1oDk/eYN0CNYC/DYV9sAv65lX8b35HYhvXISVYtwwQu/+Yg==\n-----END RSA PRIVATE KEY-----",
		3: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA02M1I1V/YvxANOvLFX8R7D8At40IlT7HDWpAW3t+tAgQ7sqj\nCeYOxiXqOaaw2kJhM3HT5nZll48UmykVq45Q05J57nhdSsGXLJshtLcTg9liMEoW\n61BjVZi9EPPRSnE05tBJc57iqZw+aEcaSU0awfzBc8IkRd6+pJ5iIgEVfuTluani\nzhHWvRli3EkAF4VNhaTfP3EkYfr4NE899aUeScbbdLFI6u1XQudlJCPTxaISx5Zc\nwM+nP3v242ABcjgUcfCbz0AY547WazK4bWP3qicyxo4MoLOoe9WBq6EuG4CuZQrz\nKnq8ltSxud/6chdg8Mqp/IasEQ2TpvY78tEXDQIDAQABAoIBAQC4uPsYk4AsSe75\n0Au6Dz7kSfIgdDhJ44AisvTmfLauMFZLtfxfjBDhCwTxuD7XnCZAxHm97Ty+AqSp\nKm/raQQsvtWalMhBqYanzjDYMRv2niJ1vGjm3WrQxBaEF+yOtvrZsK5fQTslqInI\nqknIQH7fgjazJ7Z28D18sYNj37qfFWSSymgFo+SoS/BKEr200lpRA/oaGXiHcyIO\njJidP6b7UGes7uhMXUvLrfozmCsSqslxXO5Uk5XN/fWl4LxCGX7mpNfPZIT5YBSj\nHliFkNlxIjyJg8ORLGi82M2cuyxp39r93F6uaCjLtb+rdwlGur7npgXUkKfWQJf9\nWE7uar6BAoGBAPXIuIuYFFUhqNz5CKU014jZu6Ql0z5ZA08V84cTJcfLIK4e2rqC\n8DFTldA0FtVfOGt0V08H/x2pRChGOvUwGG5nn9Dqqh6BjByUrW4z2hnXzT3ZuSDh\n6eapiCB1jl9meJ0snhF2Ps/hqWGL2b3SkCCe90qVTzOVOeLO6YUCIOq9AoGBANws\nfQkAq/0xw8neRGNTrnXimvbS+VXPIF38widljubNN7DY5cIFTQJrnTBWKbuz/t9a\nJ8QX6TFL0ci/9vhPJoThfL12vL2kWGYgWkWRPmqaBW3yz7Hs5rt+xuH3/7A5w5vm\nkEg1NZJgnsJ0rMUTu1Q6PM5CBg6OpyHY4ThBb8qRAoGAML8ciuMgtTm1yg3CPzHZ\nxZSZeJbf7K+uzlKmOBX+GkAZPS91ZiRuCvpu7hpGpQ77m6Q5ZL1LRdC6adpz+wkM\n72ix87d3AhHjfg+mzgKOsS1x0WCLLRBhWZQqIXXvRNCH/3RH7WKsVoKFG4mnJ9TJ\nLQ8aMLqoOKzSDD/JZM3lRWkCgYA8hn5Y2zZshCGufMuQApETFxhCgfzI+geLztAQ\nxHpkOEX296kxjQN+htbPUuBmGTUXcVE9NtWEF7Oz3BGocRnFrbb83odEGsmySXKH\nbUYbR/v2Ham638UOBevmcqZ3a2m6kcdYEkiH1MfP7QMRqjr1DI1qpfvERLLtOxGu\nxU5WAQKBgQCaVavyY6Slj3ZRQ7iKk9fHkge/bFl+zhANxRfWVOYMC8mD2gHDsq9C\nIdCp1Mg0tJpWLaGgyDM1kgChZYsff4jRxHC4czvAtoPSlxWXF2nL31qJ3vk2Zzzc\na4GSHAInodXBrKstav5SIKosWKT2YysxgHlA9Sm2f4n09GjFbslEHg==\n-----END RSA PRIVATE KEY-----",
		4: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAyaxqjPJP5+Innfv5IdfQqY/ftS++lnDRe3EczNkIjESWXhHS\nOljEw9b9C+/BtF+fO9QZL7Z742y06eIdvsMPQKdGflB26+9OZ8AF4SpXDn3aVWGr\n8+9qpB7BELRZI/Ph2FlFL4cobCzMHunncW8zTfMId48+fgHkAzCjRl5rC6XT0Yge\n6+eKpXmF+hr0vGYWiTzqPzTABl44WZo3rw0yurZTzkrmRE4kR2VzkjY/rBnQAbFK\nKFUKsUozjCXvSag4l461wDkhmmyivpNkK5cAxuDbsmC39iqagMt9438fajLVvYOv\npVs9ci5tiLcbBtfB4Rf/QVAkqtTm86Z0O3e7DwIDAQABAoIBAQCyma226vTW35LE\nN5zXWuAg+hhcxk6bvofWMUMXKvGF/0vHPTMXlvuSkDeDNa4vBivneRthBNPMgb3q\nDuTWxrogQMOOI8ZdhY3DFexfDvcQD2anDJuSqSmg9Nd36q+yxk3xIoXB5Ilo23dd\nvTnJXHhsBNovv7zRLO134cAHFqDoKzt5EEHre0skUcn6HjHOek6c53jvpKr5LSrr\niwx5gMuY/7ZSIUDo9WGY70qbQFGY6bOlX9x8uNjcFF+7SztEVQ+vhJ/+7EvwqaJr\nysweo0l91TKM9WaMuwoucKeceVWuynEw6GGTw8UTLtltekLGe6bS8YxY8fVwnKkT\nRwJYwAJRAoGBAP2rhcfOA+1Ja37hUHKebfp9rHsex4+pGyt3Kdu7WdqOn4sexmya\nBuiHQcUchPDVla/ruQZ20+8LHgzBDo0m8sY7gpf715UV9NSVIRD0wu26SKRklOFz\nJ4HBOwU9hBGLSnRUJzyvVlt5O7E9hAv61SCrvWBEcow2YnKNQLwvjMVJAoGBAMuG\noSb3A/ulqtp2zpxVAclYe/bSItZZTOUWP6Vb4hOiHxIJ0n1H9ap6grOYkJ/Yn4gg\nyYzKm/noF1wXP7Rj/xOahnvMkzhGdmOabvE9LH5HwQTWxBBWTkZzgBbYtbg+J5MT\ncKqJaychSRjJj+xX+d90rtlSu/c27chlSRKAHXWXAoGAFTcIHDq9l1XBmL3tRXi8\nh+uExlM/q2MgM5VmucrEbAPrke4D+Ec1drMBLCQDdkTWnPzg34qGlQJgA/8NYX61\nZSDK/j0AvaY1cKX8OvfNaaZftuf2j5ha4H4xmnGXnwQAORRkp62eUk4kUOFtLrdO\npcnXL7rpvZI6z4vCszpi0okCgYEAp3lZEl8g/+oK9UneKfYpSi1tlGTGFevVwozU\nQpWhKta1CnraogycsnOtKWvZVi9C1xljwF7YioPY9QaMfTvroY3+K9DjM+OHd96U\nfB4Chsc0pW60V10te/t+403f+oPqvLO6ehop+kEBjUwPCkQ6cQ3q8xmJYpvofoYZ\n4wdZNnECgYBwG8Vrv7Z+kX9Zuh1FvcRoY57bYLU0cWW92SA3Nvi8pZOIEaLHrQyZ\npvvaLIicR1m9+KsOAmii7ru0zL7KsrGW+5migQsaDi4gzahKQpad/R7MLKi/L53r\nYmo0aZKARLHW82GbomQ0zxdRoo9vaqfGNpXkxyyt3k3GGDunmrskYw==\n-----END RSA PRIVATE KEY-----",
		5: "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAsJbFp3WcsiojjdQtVnTuvtawL2m4XxK93F6lCnFwcZqUP39t\nxFGGlrogHMqreyawIUN7E5shtwGzigzjW8Ly5CryBJpXP3ehNTqJS7emb+9LlC19\nOxa1eQuUQnatgcsd16DPH7kJ5JzN3vXnhvUyk4Qficdmm0uk7FRaNYFi7EJs4xyq\nFTrp3rDZ0dzBHumlIeK1om7FNt6Nyivgp+UybO7kl0NLFEeSlV4S+7ofitWQsO5x\nYqKAzSzz+KIRQcxJidGBlZ1JN/g5DPDpx/ztvOWYUlM7TYk6xN3focZpU0kBzAw/\nrn94yW9z8jpXfzk+MvWzVL/HAcPy4ySwkay0NwIDAQABAoIBADzKWpawbVYEHaM4\nlLb7oCjAPXzE9zx7djLDvisfLCdfoINPedkoe52ty1o+BtRpWB7LXTY9pFic1FLE\n5wvyy6zyf8hH3ZsysqNhWFxhh4FnLmx/UGokAir+anaK5mYVJ1vQtxzjlV1HAbQs\nkRyrklKoHDdRFqiFXOwiib97oDNWhD+RxfyGwwJnynZZSXdLbLSiz/QHQNr/+Ufk\nKRBaxt0CfU7mOLZxoy6fNAxHdBcBJPHCyh+aDvEbix7nSncSU8Ju/48YJ8DrglbZ\nsXCYoA5Uz8NMDuaEMgoNWCFQVoEcRkEUoaH7BlWd3UUFRPnDZ1B4BmkrVoRE8a58\n3OqSwakCgYEA19wQUISXtpnmCrEZfbyZ6IwOy8ZCVaVUtbTjVa8UyfNglzzJG3yz\ncXU3X35v5/HNCHaXbG2qcbQLThnHBA+obW3RDo+Q49V84Zh1fUNH0ONHHuC09kB/\n/gHqzn/4nLf1aJ2O0NrMyrZNsZ0ZKUKQuVCqWjBOmTNUitcc8RpXZ8sCgYEA0W09\nPOM/It7RoVGI+cfbbgSRmzFo9kzSp5lP7iZ81bnvUMabu2nv3OeGc3Pmdh1ZJFRw\n6iDM6VVbG0uz8g+f8+JT32XdqM7MJAmgfcYfTVBMiVnh330WNkeRrGWqQzB2f2Wr\n+0vJjU8CAAcOWDh0oNguJ1l1TSyKxqdL8FsA38UCgYEAudt1AJ7psgOYmqQZ+rUl\nH6FYLAQsoWmVIk75XpE9KRUwmYdw8QXRy2LNpp9K4z7C9wKFJorWMsh+42Q2gzyo\nHHBtjEf4zPLIb8XBg3UmpKjMV73Kkiy/B4nHDr4I5YdO+iCPEy0RH4kQJFnLjEcQ\nLT9TLgxh4G7d4B2PgdjYYTkCgYEArdgiV2LETCvulBzcuYufqOn9/He9i4cl7p4j\nbathQQFBmSnkqGQ+Cn/eagQxsKaYEsJNoOxtbNu/7x6eVzeFLawYt38Vy0UuzFN5\neC54WXNotTN5fk2VnKU4VYVnGrMmCobZhpbYzoZhQKiazby/g60wUtW9u7xXzqOd\nM/428YkCgYBwbEOx1RboH8H+fP1CAiF+cqtq4Jrz9IRWPOgcDpt2Usk1rDweWrZx\nbTRlwIaVc5csIEE2X02fut/TTXr1MoXHa6s2cQrnZYq44488NsO4TAC26hqs/x/H\nbVOcX13gT26SYngAHHeh7xjWJr/KgIIwvcvgvoVs6lu7a8aLUvrOag==\n-----END RSA PRIVATE KEY-----",
	},
}
