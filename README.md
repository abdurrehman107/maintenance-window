# maintenance-window
There are times when we wish to block our DevOps team from making any deployment. We make sure we inform them, however there seems to be a need for an additional security layer using which we ensure that no deployment is made. 

The Maintenance Window CRD is meant to ensure that no deployments are allowed when the currentTime lies between startTime and endTime. 