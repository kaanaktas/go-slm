# go-slm

go-slm is a policy-based service level management library that enforces policy requirements as per service. Different policy rules can be combined 
and set to different services via basic configuration files.

Introduction
------------

Currently, go-slm supports data filtering including owasp sql injection rules, owasp xss rules and PAN process, and rule definitions for each can be found under **datafilter/rules**. 
Existing rules can be expanded according to needs, or rules that are deemed unnecessary can be disabled.
The rule-sets under https://github.com/coreruleset/coreruleset are referenced for Owasp rule definitions. 
If there is a requirement for other rule-sets in **Coreruleset**, configuration files can be created in the same way.

Installation
-------------

`go get github.com/kaanaktas/go-slm`

Usage
-------

We can create new rules in the **common_rules.json** file under the policy folder and reuse these rules 
to combine different policies in **policy_rule_set.json**.

**common_rules.json**

```
  "rules": [
    {
      "Name": "combined_rule",
      "rule": [
        {
          "name": "xss",
          "active": true
        },
        {
          "name": "pan_process",
          "active": true
        },
        {
          "name": "sqli",
          "active": true
        }
      ]
    },
    {
      "Name": "pan_only_rule",
      "rule": [
        {
          "name": "pan_process",
          "active": true
        }
      ]
    }
  ]
```

**policy_rule_set.json**

```
"policies": [
    {
        "serviceName": "test",
        "request": "combined_rule",
        "response": "pan_only_rule"
    },
    {
        "serviceName": "test2",
        "request": "combined_rule",
        "response": "combined_rule"
    }
]
```
