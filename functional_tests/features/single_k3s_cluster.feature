Feature: Kubernetes Cluster
  In order to test Kubernetes clusters
  I should apply a blueprint
  And test the output

  Scenario: K3s Cluster
    Given I apply the config "./test_fixtures/single_k3s_cluster"
    Then there should be 1 network called "cloud"
    And there should be 1 container running called "server.k3s.k8s_cluster.shipyard.run"
    And there should be 1 container running called "consul-http.ingress.shipyard.run"
    And a call to "http://consul-http.ingress.shipyard.run:18500/v1/agent/members" should result in status 200