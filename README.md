# OpenFGA Envoy

This repository contains the integration to enforce access control with OpenFGA on Envoy.

## Overview

openfga-envoy extends envoy with a gRPC server that implements the [Envoy External Authorization API](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/security/ext_authz_filter.html). You can use this integration to enforce fine-grained, context-aware access control policies with Envoy transparently to the upstream services.
