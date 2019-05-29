# Designing a URL Shortener Service

Design a url shortener service that can scale to millions.

## Requirements

Functional requirements:
- System should generate a unique short URL for every url provided.
- System should provide an option for users to choose the short name they desired, without duplication.
- System should not blacklist profanity words such as f_ck etc. We can reverse engineer and find the ids that will generate such word and ban them.
- System should allow user to check if URLs is used or not. We can use a bloom filter for better performance.
- System should generate unique/non-duplicate short url.

Non-functional requirements:
- System should be highly available. Downtime means that all the service will not redirect user to their destination.
- System should be persistent. URLs should be stored permanently. Removing a url will lead to dead links.

Extended requirements:
- System should provide analytics to users.
- System should provide login to users for them to maintain their URL list.

## Capacity and Constraints

TODO


## Use Cases

### User submits url to be shortened

Main success scenario:
1. User enters long url
2. System validates that the url exist and returns a short id

Extensions:
2a). The long url does not exist
  - System returns error
  
  
### User visits short url
Main success scenario:
1. User enters short url
2. The system validates the url exists and redirects User to the long url.

Extensions:
2a). The short url is invalid
  - System redirects user to error page


## Others

- url needs to be valid (get requests returns 200)
- short id cannot contain bad words (check the dictionary to ensure there are no bad words)
- short id cannot be predictable (using base62 will get predictable url like a, b, c….)
- how about getting the first 8 characters of the uuid v4? is this collision free? or how about auto-increment + uuid v4
- can the user choose their own short url
- how to handle expiry for short url (is this even allowed)
- can short url be reused? same short id
- limiting the number of short urls that each user can generate


| Number of characters	| Number of combinations | String Examples |
| - | - | - | 
| 1	| 62^1 = 62 | 	a, b, 1, 2| 
| 2	| 62^2 = 3.844	| 9X, aB| 
| 3	| 62^3 = 238.328	| xy4, LnK| 
| 4	| 62^4 = 14.776.336	| kY3v| 
| 5	| 62^5 = 916.132.832	| tGh32| 
| 6	| 62^6 = 56.800.235.584	| 2AjHg8| 
| 7	| 62^7 = 3.521.614.606.208	| u3gGbc8| 
| 8	| 62^8 = 218.340.105.584.896	| ShortLnk| 
| 9	| 62^9 = 13.537.086.546.263.552	| VtE4o1mng| 
| 10	| 62^10 = 839.299.365.868.340.224	| MyLink1234| 
