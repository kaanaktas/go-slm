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

# Configuration

## datafilter

Currently, go-slm implements 3 data filters, **owasp-sqli**, **owasp-xss** and **pan-filtering**. The default definitions for these filters are defined in the go-slm package 
and can be viewed under **datafilter/rules**. At the same time, the definitions of these filters are defined in **datafilter/datafilter_rule_set.json** and are ready to use without any modification.<br/>
If users want to make any changes in the existing filters, or if they want to add new rules to the filters;
* First, they need to create custom filter files and put them into the project directory. 
* Second, they need to create a custom **datafilter_rule_set.json** file and put it into the project directory. Users can update existing types/rules in the default datafilter_rule_set.json file 
or define new types/rules with changes made in their own datafilter_rule_set.json.
* Finally, custom filter files should be linked in the custom datafilter_rule_set.json. 

**custom_owasp_attack_sqli.json**

(As an example, let's assume that we put this file under the **/config** directory in the main application.)

```
[
  {
    "name": "942110",
    "disable" : true,
    "rule": "(?:^\\s*[\\\"'`;]+|[\\\"'`]+\\s*$)",
    "message": "My custom message: SQL Injection Attack: Common Injection Testing Detected",
    "sample": "var=''"
  },
  {
    "name": "new_rule_1",
    "disable" : false,
    "rule": "<new_rule_regex>",
    "message": "<new_rule_message>",
    "sample": "<new_rule_sample>"
  },
]
```

In the example file above, 2 rules are defined for owasp_attack_sqli. 
* The first rule with name=942110 updates and disables the existing rule in the package rule file (**datafilter/rules/owasp_attack_sqli.json**). 
By doing this, we disable the rule which is not required in our rule set. Similarly, we can change the rule message or regex value as needed.
* The second rule creates a new filter rule and adds it to the rule set which is generated from the package rule file.


**custom_datafilter_rule_set.json**

```
[
  {
    "type": "owasp",
    "rules": [
      {
        "name": "sqli",
        "path": "rules/owasp_attack_sqli.json"
        "custom_path": "config/custom_owasp_attack_sqli.json"
      }
    ]
  }
]
```

In the **custom_datafilter_rule_set.json** file above, we define a single rule which only updates **owasp_sqli** and leaves the other rules as is.
So, the rules inside **custom_owasp_attack_sqli.json** update the rules defined in the **owasp_attack_sqli.json** file if necessary, or add them to our rule_set as a new rule.</br>
In order for the newly created **custom_owasp_attack_sqli.json** file to be valid, it must be defined in the **GO_SLM_DATA_FILTER_RULE_SET_PATH** environment variable as in the example below.

`_ = os.Setenv("GO_SLM_DATA_FILTER_RULE_SET_PATH", "/config/custom_datafilter_rule_set.json")
`
## policy

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
