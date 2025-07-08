# Booking service

## Scope of features
- Google auth
- Types of users:
  - business accounts;
  - customers who set up appointments.
- Booking service
- Payment transaction service
- Notification service

### Example endpoints:
**/login** 
**/signup**
/account-info: CRUD operations for users data
- as a user I can manage my info

**/my-bookings** 
- as a user I can see the list of my bookings  

**/business-accounts**: CRUD operations for managing business accounts.
- as a user I can create a business account to provide my services
- as a user who has a business account, I can manage a business account to provide my services

**/services**: CRUD operations for managing services.
- as a user who has a business account, I can configure the services that I provide:

**/schedules**: Endpoint to retrieve available time slots for a given master and service.
- as a user, I can see the list of available time slots  for the selected service of some business account

**/bookings**: Endpoint to handle bookings.
- as a user, I can book some service of some business account user
- as a user, I can reschedule my booking or cancel it but not less than 1 day
- a business account user should be notified about booking changes 
- as a user who has a business account, I can manually add some appointment to my schedule
- as a user who has a business account, I can reschedule my booking or cancel it
- a user should be notified if business account updated the booking.

**/notifications**
- as a user I want to see the list of notifications.
