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
