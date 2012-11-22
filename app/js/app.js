/*
Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the 
terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
*/

function DecisionCtrl($scope, $http) {
  $scope.master = {
		  quandary: "",
		  choices:  [
		             {"text": ""},
		             {"text": ""},
		             {"text": ""}
		             ]
  };
  
  $scope.orig = angular.copy($scope.master);
  
 
  $scope.update = function(decision) {
    $scope.master= angular.copy(decision);
  };
 
  $scope.reset = function() {
    $scope.decision = angular.copy($scope.master);
    $scope.winner = null;
	$scope.error = null;
  };
 
  $scope.reset();
  
  $scope.winner = null;
  
  $scope.decide = function(decision) {
	  $http.post("/v1/decide", decision).
	  	success(function(data, status){
	  		$scope.error = null;
	  		$scope.data = data;
	  		$scope.winner = data["Decision"];
	  	}).
	  error(function(data, status) {
		  	$scope.winner = null;
	  		$scope.error = data || "Request Failed";
	  });
  };
}