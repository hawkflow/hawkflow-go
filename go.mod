module github.com/hawkflow/hawkflow-go

go 1.17

retract (
    v1.0.0 // Contains insufficiently small timeout.
    v1.0.3 // Contains invalid regex validation.
)