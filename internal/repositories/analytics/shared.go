package analytics

import "git.mos.ru/buch-cloud/moscow-team-2.0/build/diterrors.git"

const (
    ErrMetricsNotProvided = diterrors.StringError("metrics not provided")
    ErrEmptyResponse = diterrors.StringError("empty response")
    ErrEmptyRequest = diterrors.StringError("empty request")
)
