import { Command } from 'commander';
import inquirer from 'inquirer';
import axios from 'axios';

const program = new Command();

const USERS_BASE_URL = process.env.USERS_SERVICE_URL || 'http://localhost:3000';
const RIDES_BASE_URL = process.env.RIDES_SERVICE_URL || 'http://localhost:8080';

interface Driver {
  id: string;
  name: string;
  is_available: boolean;
}

interface Passenger {
  id: string;
  name: string;
  created_at: string;
  updated_at: string;
}

interface Ride {
  id: string;
  passengerId: string;
  driverId: string;
  from_zone: string;
  to_zone: string;
  price: number;
  status: 'ASSIGNED' | 'IN_PROGRESS' | 'COMPLETED' | 'CANCELLED';
  paymentStatus: 'PENDING' | 'CAPTURED';
  createdAt: string;
  updatedAt: string;
}

const usersApi = axios.create({
  baseURL: USERS_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

const ridesApi = axios.create({
  baseURL: RIDES_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

const handleError = (error: any) => {
  if (axios.isAxiosError(error)) {
    console.error('Error:', error.response?.data || error.message);
  } else {
    console.error('Error:', error.message);
  }
};

const displayIntro = () => {
  console.log('========================================');
  console.log('   MGL7361-Microservices implementation');
  console.log('   Users & Rides Service CLI');
  console.log('========================================');
  console.log('');
};

const setupDriverCommands = (program: Command) => {
  const driverCommand = program.command('drivers');
  
  driverCommand
    .command('create')
    .description('Create a new driver')
    .action(async () => {
      try {
        const answers = await inquirer.prompt([
          {
            type: 'input',
            name: 'name',
            message: 'Enter driver name:',
            validate: (input) => input.trim() ? true : 'Name is required'
          }
        ]);

        const response = await usersApi.post('/drivers', {
          name: answers.name
        });

        console.log('Driver created successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  driverCommand
    .command('list')
    .description('List all drivers')
    .option('-a, --available', 'Show only available drivers')
    .action(async (options) => {
      try {
        const params = options.available ? { available: true } : {};
        const response = await usersApi.get('/drivers', { params });
        
        console.log('Drivers:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  driverCommand
    .command('update-status')
    .description('Update driver availability status')
    .action(async () => {
      try {
        const driversResponse = await usersApi.get('/drivers');
        const drivers: Driver[] = driversResponse.data;

        const driverAnswer = await inquirer.prompt([
          {
            type: 'list',
            name: 'driverId',
            message: 'Select driver:',
            choices: drivers.map(driver => ({
              name: `${driver.name} (${driver.is_available ? 'Available' : 'Unavailable'})`,
              value: driver.id
            }))
          }
        ]);

        const statusAnswer = await inquirer.prompt([
          {
            type: 'confirm',
            name: 'isAvailable',
            message: 'Set as available?',
            default: true
          }
        ]);

        const response = await usersApi.patch(`/drivers/${driverAnswer.driverId}/status`, {
          is_available: statusAnswer.isAvailable
        });

        console.log('Driver status updated successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });
};

const setupPassengerCommands = (program: Command) => {
  const passengerCommand = program.command('passengers');
  
  passengerCommand
    .command('create')
    .description('Create a new passenger')
    .action(async () => {
      try {
        const answers = await inquirer.prompt([
          {
            type: 'input',
            name: 'name',
            message: 'Enter passenger name:',
            validate: (input) => input.trim() ? true : 'Name is required'
          }
        ]);

        const response = await usersApi.post('/passengers', {
          name: answers.name
        });

        console.log('Passenger created successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('list')
    .description('List all passengers')
    .action(async () => {
      try {
        const response = await usersApi.get('/passengers');
        
        console.log('Passengers:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('get')
    .description('Get passenger by ID')
    .action(async () => {
      try {
        const passengersResponse = await usersApi.get('/passengers');
        const passengers: Passenger[] = passengersResponse.data;

        const answers = await inquirer.prompt([
          {
            type: 'list',
            name: 'passengerId',
            message: 'Select passenger:',
            choices: passengers.map(passenger => ({
              name: `${passenger.name} (${passenger.id})`,
              value: passenger.id
            }))
          }
        ]);

        const response = await usersApi.get(`/passengers/${answers.passengerId}`);
        console.log('Passenger details:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('update')
    .description('Update passenger information')
    .action(async () => {
      try {
        const passengersResponse = await usersApi.get('/passengers');
        const passengers: Passenger[] = passengersResponse.data;

        const passengerAnswer = await inquirer.prompt([
          {
            type: 'list',
            name: 'passengerId',
            message: 'Select passenger to update:',
            choices: passengers.map(passenger => ({
              name: `${passenger.name} (${passenger.id})`,
              value: passenger.id
            }))
          }
        ]);

        const nameAnswer = await inquirer.prompt([
          {
            type: 'input',
            name: 'name',
            message: 'Enter new name:',
            validate: (input) => input.trim() ? true : 'Name is required'
          }
        ]);

        const response = await usersApi.put(`/passengers/${passengerAnswer.passengerId}`, {
          name: nameAnswer.name
        });

        console.log('Passenger updated successfully:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  passengerCommand
    .command('delete')
    .description('Delete a passenger')
    .action(async () => {
      try {
        const passengersResponse = await usersApi.get('/passengers');
        const passengers: Passenger[] = passengersResponse.data;

        const answers = await inquirer.prompt([
          {
            type: 'list',
            name: 'passengerId',
            message: 'Select passenger to delete:',
            choices: passengers.map(passenger => ({
              name: `${passenger.name} (${passenger.id})`,
              value: passenger.id
            }))
          },
          {
            type: 'confirm',
            name: 'confirm',
            message: 'Are you sure you want to delete this passenger?',
            default: false
          }
        ]);

        if (answers.confirm) {
          await usersApi.delete(`/passengers/${answers.passengerId}`);
          console.log('Passenger deleted successfully');
        } else {
          console.log('Deletion cancelled');
        }
      } catch (error) {
        handleError(error);
      }
    });
};

const setupRideCommands = (program: Command) => {
  const rideCommand = program.command('rides');
  
  rideCommand
    .command('create')
    .description('Create a new ride (reserve a ride)')
    .action(async () => {
      try {
        // Get all passengers
        const passengersResponse = await usersApi.get('/passengers');
        const passengers: Passenger[] = passengersResponse.data;

        if (passengers.length === 0) {
          console.log('No passengers found. Please create a passenger first.');
          return;
        }

        const answers = await inquirer.prompt([
          {
            type: 'list',
            name: 'passengerId',
            message: 'Select passenger:',
            choices: passengers.map(passenger => ({
              name: `${passenger.name} (${passenger.id})`,
              value: passenger.id
            }))
          },
          {
            type: 'input',
            name: 'from_zone',
            message: 'Enter pickup zone:',
            validate: (input) => input.trim() ? true : 'Pickup zone is required'
          },
          {
            type: 'input',
            name: 'to_zone',
            message: 'Enter destination zone:',
            validate: (input) => input.trim() ? true : 'Destination zone is required'
          }
        ]);

        const response = await ridesApi.post('/rides', {
          passengerId: answers.passengerId,
          from_zone: answers.from_zone,
          to_zone: answers.to_zone
        });

        console.log('Ride created successfully:');
        console.log(JSON.stringify(response.data, null, 2));
        
        const ride: Ride = response.data;
        console.log(`\nSummary:`);
        console.log(`- Passenger: ${answers.passengerId}`);
        console.log(`- From: ${ride.from_zone}`);
        console.log(`- To: ${ride.to_zone}`);
        console.log(`- Driver assigned: ${ride.driverId}`);
        console.log(`- Price: $${ride.price}`);
        console.log(`- Status: ${ride.status}`);
        console.log(`- Payment Status: ${ride.paymentStatus}`);
      } catch (error) {
        handleError(error);
      }
    });

  rideCommand
    .command('get')
    .description('Get ride details by ID')
    .action(async () => {
      try {
        const ridesResponse = await ridesApi.get('/rides');
        const rides: Ride[] = ridesResponse.data;

        if (rides.length === 0) {
          console.log('No rides found.');
          return;
        }

        const answers = await inquirer.prompt([
          {
            type: 'list',
            name: 'rideId',
            message: 'Select ride:',
            choices: rides.map(ride => ({
              name: `Ride ${ride.id} - ${ride.from_zone} to ${ride.to_zone} (${ride.status})`,
              value: ride.id
            }))
          }
        ]);

        const response = await ridesApi.get(`/rides/${answers.rideId}`);
        console.log('Ride details:');
        console.log(JSON.stringify(response.data, null, 2));
      } catch (error) {
        handleError(error);
      }
    });

  rideCommand
    .command('list')
    .description('List all rides')
    .option('-s, --status <status>', 'Filter by status (ASSIGNED, IN_PROGRESS, COMPLETED, CANCELLED)')
    .action(async (options) => {
      try {
        const params = options.status ? { status: options.status } : {};
        const response = await ridesApi.get('/rides', { params });
        
        if (response.data.length === 0) {
          console.log('No rides found.');
          return;
        }
        
        console.log('Rides:');
        response.data.forEach((ride: Ride) => {
          console.log(`\nID: ${ride.id}`);
          console.log(`Passenger: ${ride.passengerId}`);
          console.log(`Driver: ${ride.driverId}`);
          console.log(`From: ${ride.from_zone}`);
          console.log(`To: ${ride.to_zone}`);
          console.log(`Price: $${ride.price}`);
          console.log(`Status: ${ride.status}`);
          console.log(`Payment: ${ride.paymentStatus}`);
          console.log(`Created: ${new Date(ride.createdAt).toLocaleString()}`);
          console.log('---');
        });
      } catch (error) {
        handleError(error);
      }
    });

  rideCommand
    .command('update-status')
    .description('Update ride status')
    .action(async () => {
      try {
        const ridesResponse = await ridesApi.get('/rides');
        const rides: Ride[] = ridesResponse.data;

        if (rides.length === 0) {
          console.log('No rides found.');
          return;
        }

        const rideAnswer = await inquirer.prompt([
          {
            type: 'list',
            name: 'rideId',
            message: 'Select ride to update:',
            choices: rides.map(ride => ({
              name: `Ride ${ride.id} - ${ride.from_zone} to ${ride.to_zone} (${ride.status})`,
              value: ride.id
            }))
          },
          {
            type: 'list',
            name: 'status',
            message: 'Select new status:',
            choices: [
              { name: 'ASSIGNED - Ride assigned to driver', value: 'ASSIGNED' },
              { name: 'IN_PROGRESS - Ride in progress', value: 'IN_PROGRESS' },
              { name: 'COMPLETED - Ride completed', value: 'COMPLETED' },
              { name: 'CANCELLED - Ride cancelled', value: 'CANCELLED' }
            ]
          }
        ]);

        const response = await ridesApi.patch(`/rides/${rideAnswer.rideId}/status`, {
          status: rideAnswer.status
        });

        console.log('Ride status updated successfully:');
        console.log(JSON.stringify(response.data, null, 2));
        
        const ride: Ride = response.data;
        console.log(`\nUpdated Status: ${ride.status}`);
        console.log(`Payment Status: ${ride.paymentStatus}`);
        
        if (ride.status === 'COMPLETED') {
          console.log('\nNote: The ride has been marked as completed.');
          console.log('- Payment has been automatically captured');
          console.log('- Driver has been marked as available again');
        }
      } catch (error) {
        handleError(error);
      }
    });

  rideCommand
    .command('status')
    .description('Get current ride status and details')
    .action(async () => {
      try {
        const ridesResponse = await ridesApi.get('/rides');
        const rides: Ride[] = ridesResponse.data;

        if (rides.length === 0) {
          console.log('No rides found.');
          return;
        }

        const answers = await inquirer.prompt([
          {
            type: 'list',
            name: 'rideId',
            message: 'Select ride:',
            choices: rides.map(ride => ({
              name: `Ride ${ride.id} - ${ride.from_zone} to ${ride.to_zone}`,
              value: ride.id
            }))
          }
        ]);

        const response = await ridesApi.get(`/rides/${answers.rideId}`);
        const ride: Ride = response.data;

        console.log('\n=== Ride Status ===');
        console.log(`Ride ID: ${ride.id}`);
        console.log(`Passenger ID: ${ride.passengerId}`);
        console.log(`Driver ID: ${ride.driverId}`);
        console.log(`Route: ${ride.from_zone} → ${ride.to_zone}`);
        console.log(`Price: $${ride.price}`);
        console.log(`Status: ${ride.status}`);
        console.log(`Payment Status: ${ride.paymentStatus}`);
        console.log(`Created: ${new Date(ride.createdAt).toLocaleString()}`);
        console.log(`Last Updated: ${new Date(ride.updatedAt).toLocaleString()}`);
        
        // Get driver details
        try {
          const driverResponse = await usersApi.get(`/drivers/${ride.driverId}`);
          console.log(`\nDriver: ${driverResponse.data.name}`);
          console.log(`Driver Available: ${driverResponse.data.is_available ? 'Yes' : 'No'}`);
        } catch (error) {
          console.log('\nDriver details not available');
        }
      } catch (error) {
        handleError(error);
      }
    });
};

const main = async () => {
  displayIntro();

  program
    .name('ridenow-cli')
    .description('CLI for Users & Rides Service API')
    .version('1.0.0');

  setupDriverCommands(program);
  setupPassengerCommands(program);
  setupRideCommands(program);

  program
    .command('interactive')
    .description('Start interactive mode')
    .action(async () => {
      while (true) {
        const { action } = await inquirer.prompt([
          {
            type: 'list',
            name: 'action',
            message: 'What would you like to do?',
            choices: [
              { name: 'Manage Drivers', value: 'drivers' },
              { name: 'Manage Passengers', value: 'passengers' },
              { name: 'Manage Rides', value: 'rides' },
              { name: 'Exit', value: 'exit' }
            ]
          }
        ]);

        if (action === 'exit') {
          console.log('Goodbye!');
          process.exit(0);
        }

        if (action === 'drivers') {
          const { driverAction } = await inquirer.prompt([
            {
              type: 'list',
              name: 'driverAction',
              message: 'Driver Management:',
              choices: [
                { name: 'Create Driver', value: 'create' },
                { name: 'List Drivers', value: 'list' },
                { name: 'Update Driver Status', value: 'update-status' },
                { name: 'Back to Main Menu', value: 'back' }
              ]
            }
          ]);

          if (driverAction === 'back') continue;

          const driverCommand = program.commands.find(cmd => cmd.name() === 'drivers');
          if (driverCommand) {
            const subCommand = driverCommand.commands.find(cmd => cmd.name() === driverAction);
            if (subCommand) {
              await subCommand.parseAsync([], { from: 'user' });
            }
          }
        }

        if (action === 'passengers') {
          const { passengerAction } = await inquirer.prompt([
            {
              type: 'list',
              name: 'passengerAction',
              message: 'Passenger Management:',
              choices: [
                { name: 'Create Passenger', value: 'create' },
                { name: 'List Passengers', value: 'list' },
                { name: 'Get Passenger Details', value: 'get' },
                { name: 'Update Passenger', value: 'update' },
                { name: 'Delete Passenger', value: 'delete' },
                { name: 'Back to Main Menu', value: 'back' }
              ]
            }
          ]);

          if (passengerAction === 'back') continue;

          const passengerCommand = program.commands.find(cmd => cmd.name() === 'passengers');
          if (passengerCommand) {
            const subCommand = passengerCommand.commands.find(cmd => cmd.name() === passengerAction);
            if (subCommand) {
              await subCommand.parseAsync([], { from: 'user' });
            }
          }
        }

        if (action === 'rides') {
          const { rideAction } = await inquirer.prompt([
            {
              type: 'list',
              name: 'rideAction',
              message: 'Ride Management:',
              choices: [
                { name: 'Create Ride (Reserve)', value: 'create' },
                { name: 'Get Ride Details', value: 'get' },
                { name: 'List All Rides', value: 'list' },
                { name: 'Update Ride Status', value: 'update-status' },
                { name: 'Check Ride Status', value: 'status' },
                { name: 'Back to Main Menu', value: 'back' }
              ]
            }
          ]);

          if (rideAction === 'back') continue;

          const rideCommand = program.commands.find(cmd => cmd.name() === 'rides');
          if (rideCommand) {
            const subCommand = rideCommand.commands.find(cmd => cmd.name() === rideAction);
            if (subCommand) {
              await subCommand.parseAsync([], { from: 'user' });
            }
          }
        }
      }
    });

  // Mode Interactif
  const isDocker = process.env.IS_DOCKER || process.cwd() === '/app';
  if (isDocker || process.argv.length <= 2) {
    console.log('Starting interactive mode...');
    program.commands.find(cmd => cmd.name() === 'interactive')?.parseAsync();
  } else {
    program.parse();
  }
};

main().catch(console.error);