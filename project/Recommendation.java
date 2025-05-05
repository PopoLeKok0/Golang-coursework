public class Recommendation {
    private final User user;  
    private final Movie movie;
    private final double probability;
    private final int nUsers;

    // Constructor
    public Recommendation(User user, Movie movie, double probability, int nUsers) {
        this.user = user; 
        this.movie = movie;
        this.probability = probability;
        this.nUsers = nUsers;
    }

    // gets the user
    public User getUser() {
        return user;
    }

    // gets the movie
    public Movie getMovie() { 
        return movie; 
    }

    // gets the probability
    public double getProbability() { 
        return probability; 
    }

    // gets the number of users
    public int getNUsers() { 
        return nUsers; 
    }
}
