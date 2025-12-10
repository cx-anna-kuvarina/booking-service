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
  - Create, read, update, delete services
  - Set service pricing, duration, and categories
  - Manage service availability
  - List services by business account

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

## New Features Added

### Services Management System ✅
- **Complete CRUD operations** for business services
- **Service configuration** including pricing, duration, and categories
- **Business account integration** for service ownership
- **Comprehensive validation** and error handling
- **Pagination support** for service listings
- **Filtering capabilities** by business account, category, and status

#### Service Endpoints:
- `POST /api/services/` - Create a new service
- `GET /api/services/{id}` - Get service details
- `PUT /api/services/{id}` - Update service
- `DELETE /api/services/{id}` - Delete service
- `GET /api/services/` - List services with filtering and pagination
- `GET /api/services/business-account/{id}` - Get services by business account

#### Service Features:
- Service name, description, and category
- Duration in minutes
- Pricing with currency support
- Active/inactive status management
- Automatic timestamp tracking
- Business account ownership validation

## Database Schema

The service includes the following core tables:
- `users` - User account information
- `business_accounts` - Business account details
- `services` - Service offerings with pricing and scheduling
- `bookings` - Appointment bookings
- `user_business_accounts` - User-business account relationships

## Getting Started

1. **Database Setup**: Run the Liquibase migrations to create the required tables
2. **Configuration**: Set up your environment variables for database connections and JWT secrets
3. **Authentication**: Use Google OAuth for user authentication
4. **API Usage**: All endpoints require JWT authentication via the Authorization header

## Documentation

- [Services API Documentation](docs/services-api.md) - Complete API reference for service management
- Database migrations are located in the `db/` directory
- API handlers and business logic in `internal/api/rest/`
- Data access layer in `internal/store/`

## Next Steps

The following features are planned for future development:
- Schedule management system
- Payment processing integration
- Notification system
- Calendar integration
- Review and rating system
- Advanced search and discovery
