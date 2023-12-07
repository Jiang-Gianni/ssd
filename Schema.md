{{define "main"}}

```d2
explanation: |md
  # LLMs
  The Large Language Model (LLM) is a powerful AI\
    system that learns from vast amounts of text data.\
  By analyzing patterns and structures in language,\
  it gains an understanding of grammar, facts,\
  and even some reasoning abilities. As users input text,\
  the LLM predicts the most likely next words or phrases\
  to create coherent responses. The model\
  continuously fine-tunes its output, considering both the\
  user's input and its own vast knowledge base.\
  This cutting-edge technology enables LLM to generate human-like text,\
  making it a valuable tool for various applications.
| {
  near: center-left
}

ML Platform -> Pre-trained models
ML Platform -> Model registry
ML Platform -> Compiler
ML Platform -> Validation
ML Platform -> Auditing

Model registry -> Server.Batch Predictor
Server.Online Model Server
```

```d2
{{range $tableName, $table := .TableMap}}
{{if regex "^SRP" $tableName}}
{{$table.Name}}: {
  shape: sql_table
  link: "#SRP_T_ATTACHMENT"
  {{range $columnName := $table.Columns}}{{$columnName}}: {{$col := index $table.ColumnMap $columnName}}{{$col.Type}}
  {{end}}
}
{{end}}
{{end}}
```



## End



{{end}}
