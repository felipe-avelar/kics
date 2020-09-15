package Cx

#CxPragma "$.resource.aws_ecr_repository_policy"

#Public ECR policy
#https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ecr_repository_policy

result [ getMetadata({"id" : input.All[i].CxId, "data" : [pol], "search": concat("+", ["aws_ecr_repository_policy", name])}) ] {
	pol := input.All[i].resource.aws_ecr_repository_policy[name].policy
    re_match("\"Principal\"\\s*:\\s*\"*\"", pol)
}


getMetadata(id) = res {
	some cnt
    input.All[cnt].CxId = id.id
    res := {
        "id" : input.All[cnt].CxId,
        "file" : input.All[cnt].CxFile,
        "name" : "Public ECR policy",
        "severity": "Medium",
        "cnt" : cnt,
        "search": id.search,
        "data" : id.data
    }
}