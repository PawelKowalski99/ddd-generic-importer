# Customer importer

## What?
Customer Importer is a package that concurrently parses csv file and counts domain customer into sorted slice.

## Why?
Show how many customers are using the exact domain.

## Why Worker Pool not fan-in fan-out?
For this case much safer solution is worker pool which allows client to manage amount of workers that are created at once. If you want to have similar result as FIFO try setting worker amount at least equal to customer lines in csv file (NOT RECOMMENDED)

## What is external and internal Verifier?

Internal Verifier is our own logic to verify if mail is valid without calling any external dependencies or dns/smtp servers.

External Verifier is logic that calls if dns/smtp servers are present. If the domain is valid. It makes use of the open source library ```github.com/AfterShip/email-verifier```
 but it is also possible to make our own solution such as the internalVerifier.


Recommended amount of workers for the Internal and External Verifier:

ExternalVerifier - The worker amount should be not high due to network calling. The smtp/dns providers can block.

InternalVerifier - Due to the fact that this verifier is not dependent on external dependencies such as calling servers etc the workerAmount can be set up to the limits of the machine. But it is not beneficial if you provide more workers than there are lines in file.

## Services
There are three services

Importer - responsible for importing a to z -> calling counter to get data and calling writer to save data
Writer - responsible for saving data
Counter - responsible for counting given aggregate






## Conclusions after manual testing...

If verifier is not dependent on any resourceful dependencies then even 1 worker/sequential can make the job to run in 5sec for 1.5kk Lines of csv

If verifier is somehow dependent: that is externalVerifier is used or internalVerifier somehows slows down(sleep) then every worker influences the overall result and finish time. 

If I set the InternalVerifier to sleep randomly 1-20nano seconds then for 1.5kk lines 300 workers did the work in 44secs. 400 workers did that in 5ecs and adding more workers was not beneficial. Basically if set for 400+ workers then the worker pool always has got some worker free. 