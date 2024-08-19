# gogat-api
Api gateway em go

Com o que a Api Gateway lida?

## Roteamento de Requisições (Request Routing)

Descrição: A API Gateway deve ser capaz de rotear as requisições para o serviço backend correto com base na URL, método HTTP, cabeçalhos, ou outros parâmetros.
Funcionamento na Prática: Quando uma requisição chega à API Gateway, ela analisa a rota solicitada e, com base em uma tabela de roteamento configurada, encaminha a requisição para o serviço apropriado. Por exemplo, uma requisição GET /users pode ser roteada para o serviço de usuários.

```routes:
  - path: /users
    method: GET
    service: user-service
    backend_url: http://user-service.local/api/users
```

## Autenticação e Autorização (Authentication and Authorization)
Descrição: A API Gateway deve validar a identidade do cliente (Autenticação) e verificar se ele tem permissão para acessar o recurso solicitado (Autorização).
Funcionamento na Prática: A autenticação pode ser feita usando tokens JWT, onde a API Gateway valida o token, verifica sua validade e extrai as informações de identidade. A autorização pode ser baseada em papéis (RBAC - Role-Based Access Control) ou escopos definidos no token.

```authentication:
  type: jwt
  public_key: /path/to/public.key
  algorithms: [RS256]

authorization:
  roles:
    admin:
      - /admin/*
    user:
      - /users/*
```

## Transformação de Requisições e Respostas (Request and Response Transformation)
Descrição: Permite modificar o conteúdo das requisições e respostas, como cabeçalhos, parâmetros de URL, ou o corpo da mensagem.
Funcionamento na Prática: A API Gateway pode adicionar ou remover cabeçalhos, alterar o corpo da requisição para se adaptar ao formato esperado pelo backend, ou transformar a resposta antes de enviá-la de volta ao cliente.

```transformations:
  requests:
    - add_header: 
        Authorization: "Bearer some-static-token"
    - modify_body:
        from: "application/xml"
        to: "application/json"
  responses:
    - modify_status: 
        from: 200
        to: 201
```

## Rate Limiting e Throttling
Descrição: Controla o número de requisições que um cliente pode fazer dentro de um período de tempo, prevenindo abusos e protegendo os serviços backend.
Funcionamento na Prática: Rate limiting pode ser configurado por usuário, IP, ou chave de API. A API Gateway rejeita ou retarda requisições que excedem o limite configurado, geralmente retornando um erro 429 Too Many Requests.

```rate_limiting:
  enabled: true
  limit: 1000
  period: 60s
  per_client: true
```

## Balanceamento de Carga (Load Balancing)
Descrição: A API Gateway deve distribuir as requisições de forma equilibrada entre múltiplas instâncias de serviço backend para garantir alta disponibilidade e escalabilidade.
Funcionamento na Prática: A distribuição pode ser feita com diferentes algoritmos, como round-robin, least connections, ou random. A Gateway mantém o estado dos serviços e ajusta o roteamento dinamicamente.

```load_balancing:
  strategy: round-robin
  services:
    - user-service-instance-1
    - user-service-instance-2
    - user-service-instance-3
```

## Cache de Respostas (Response Caching)
Descrição: Armazena em cache as respostas de certos endpoints para reduzir a carga nos serviços backend e melhorar o tempo de resposta.
Funcionamento na Prática: Quando uma requisição é feita para um endpoint com cache habilitado, a API Gateway verifica se a resposta está no cache e a retorna, caso contrário, encaminha a requisição para o backend e armazena a resposta para futuras requisições.

```
caching:
  enabled: true
  ttl: 60s
  paths:
    - /users
    - /products
```

## Monitoramento e Logging
Descrição: A API Gateway deve coletar métricas e logs detalhados de todas as requisições e respostas, permitindo análise e diagnóstico de problemas.
Funcionamento na Prática: Logs devem incluir informações sobre as requisições (rota, método, status HTTP, tempo de resposta) e erros. Métricas podem ser expostas para sistemas de monitoramento como Prometheus.

```logging:
  enabled: true
  level: info
  format: json
  destinations:
    - file: /var/log/api-gateway.log
    - syslog: localhost:514

metrics:
  enabled: true
  prometheus_endpoint: /metrics
```

## SSL/TLS Termination
Descrição: A API Gateway deve gerenciar as conexões seguras (HTTPS), terminando o SSL/TLS na fronteira da rede.
Funcionamento na Prática: A API Gateway termina a conexão SSL/TLS e encaminha o tráfego não criptografado para os serviços backend, ou pode recriptografá-lo antes de repassá-lo. Também pode renovar automaticamente os certificados com serviços como Let's Encrypt.

```ssl:
  enabled: true
  cert_file: /path/to/cert.pem
  key_file: /path/to/key.pem
  ca_certs:
    - /path/to/ca.pem
```

## Circuit Breaker
Descrição: Implementa o padrão de Circuit Breaker para prevenir que falhas em um serviço backend provoquem falhas cascatas na arquitetura.
Funcionamento na Prática: Quando um serviço backend começa a falhar, o Circuit Breaker abre, e a API Gateway rejeita imediatamente as requisições para aquele serviço, retornando uma resposta de fallback ou erro customizado.

```circuit_breaker:
  enabled: true
  failure_threshold: 50
  timeout: 60s
  fallback_response: "Service temporarily unavailable, please try again later."
```

## Service Discovery

```service_discovery:
  enabled: true
  registry: consul
  consul:
    address: localhost:8500
    service_tag: api-backend

```
