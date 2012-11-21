function Decision($scope) {
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
}