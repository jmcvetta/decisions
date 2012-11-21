function DecisionCtrl($scope, $http) {
  $scope.master = {
		  quandary: "",
		  choices:  [
		             {"text": ""},
		             {"text": ""},
		             {"text": ""}
		             ]
  };
 
  $scope.update = function(decision) {
    $scope.master= angular.copy(decision);
  };
 
  $scope.reset = function() {
    $scope.decision = angular.copy($scope.master);
  };
 
  $scope.reset();
  
  $scope.winner = null;
  
  $scope.decide = function(decision) {
	  $scope.master= angular.copy(decision);
	  $scope.code = null;
	  $scope.data = null;
	  $http.post("/v1/decide", $scope.master).
	  	success(function(data, status){
	  		$scope.data = data;
	  		$scope.winner = data["Decision"];
	  		$scope.status = status;
	  	}).
	  error(function(data, status) {
	  		$scope.data = data || "Request Failed";
	  		$scope.status = status;
	  });
  };
}