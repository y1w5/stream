// This program implements a streaming pipeline between a database and a controller.
using System;

// User represents an user in our awesome banking app.
public class User {
    public string name;
    public int balance;

    public User(string name, int balance) {
        this.name = name;
        this.balance = balance;
    }
}

// DB represents the database access layer of our app.
public static class DB {

    // users mocks our database.
    public static User[] users = {
        new User("Flamboyant Greider", 42238),
        new User("Xenodochial Kirch", 2780),
        new User("Nostalgic Borg", 40312),
        new User("Epic Brahmagupta", 51939),
        new User("Fervent Keller", 55383),
        new User("Cool Bose", 59103),
        new User("Elegant Agnesi", 5395),
        new User("Quirky Joliot", 42001),
        new User("Pedantic Dijkstra", 27688),
        new User("Intelligent Kilby", 41285),
        new User("Relaxed Boyd", 32100),
        new User("Boring Buck", 26354),
        new User("Interesting Blackwell", 57291),
        new User("Elastic Faraday", 8497),
        new User("Kind Ritchie", 52535),
        new User("Determined Euclid", 34358),
        new User("Clever Jennings", 17821),
        new User("Keen Wozniak", 24441),
        new User("Silly Napier", 23884),
        new User("Vigorous Brahmagupta", 23478),
    };

    // streamUsers streams users from the database.
    //
    // TODO: fn should take an array of User for improved performance.
    public static void streamUsers(Action<User> fn) {
        foreach (User u in DB.users)
        {
            fn(u);
        }
    }
}

// Service stores the business logic of our application.
public static class Service {

    // streamUsers streams users from the database.
    public static void streamUsers(Action<User> fn) {
        DB.streamUsers((u) => {
                // We remove the poor, they are not paying enough for our services.
                if (u.balance < 20000) {
                    return;
                }
                fn(u);
        });
    }

}

public class Program
{
	public static void Main()
	{
        // We use Main as our controller.
        // The following code should be in the HTTP handler of a web service.
        
        // We write a fake HTTP header. For a valid request use chunked encoding.
        Console.Write(@"HTTP/1.1 200 OK
Content-Type: application/json
Transfer-Encoding: chunked

");

        Console.Write("[");
        var count = 0;
        Service.streamUsers((u) => {
            var prefix = "\n\t";
            if (count > 0) {
                prefix = ",\n\t";
            }

            Console.Write($"{prefix}{{\"name\": \"{u.name}\", \"balance\": \"{u.balance}\"}}");
            System.Threading.Thread.Sleep(250);
            count++;
        });
        Console.Write("\n]\n\n");
	}
}
