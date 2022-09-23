import json
import random
import pickle
import sys
import os
from sklearn.svm import SVC
from sklearn.feature_extraction.text import CountVectorizer
from sklearn.model_selection import train_test_split
from sklearn.metrics import f1_score


def train():
    # get directory of the current file
    with open("/messages.json") as f:
        file = json.load(f)

    messages = process_data(file)

    training, testing = train_test_split(messages, test_size=0.33, shuffle=True)

    train_container = TestContainer(training)
    test_container = TestContainer(testing)

    train_x = train_container.get_messages()  # the data we give to your model
    train_y = train_container.get_is_ivan()  # the labels we want to predict

    test_x = test_container.get_messages()  # the data we give to your model
    test_y = test_container.get_is_ivan()  # the labels we want to predict

    vectorizer = CountVectorizer()

    train_x_vector = vectorizer.fit_transform(train_x)
    test_x_vector = vectorizer.transform(test_x)

    # svm = SVC(kernel='linear')
    svm = SVC(kernel="rbf")
    return svm, test_x_vector, test_y, train_x_vector, train_y


class Message:
    """
    A message object
    """

    message = ""
    is_ivan = False

    def __init__(self, message, is_ivan):
        self.message = message
        self.is_ivan = is_ivan

    def __str__(self):
        return "Message: " + self.message + " is_ivan: " + str(self.is_ivan)


class TestContainer:
    """
    Container for the test data.
    """

    def __init__(self, messages):
        """
        Initializes the test container.

        Args:
            messages (): List of messages.
        """
        self.messages = messages

    def get_messages(self):
        """
        Returns a list of messages.

        Returns:
            list[str]: List of messages.

        """
        return [x.message for x in self.messages]

    def get_is_ivan(self):
        """
        Returns a list of booleans indicating whether the message is from Ivan or not.

        Returns:
            list[bool]: List of booleans indicating whether the message is from Ivan or not.
        """
        return [x.is_ivan for x in self.messages]


def process_data(messages):
    """
    Processes the data and returns a list of tuple of Messages.
    Each tuple contains the message and if it's ivan or not.

    Args:
        messages: json object containing the messages.

    Returns:
        list of tuples of shuffled Message objects.
    """
    data = []
    for message in messages:
        data.append(Message(message["message"], message["is_ivan"]))
    # shuffle the data
    random.shuffle(data)
    return data


def save_model(model, filename):
    """
    Saves the model to a file.

    Args:
        model: The model to save.
        filename: The filename to save the model to.
    """
    pickle.dump(model, open(filename, "wb"))


def load_model(filename):
    """
    Loads the model from a file.

    Args:
        filename: The filename to load the model from.

    Returns:
        The model.
    """
    return pickle.load(open(filename, "rb"))


def main():
    with open("./messages.json") as f:
        file = json.load(f)

    messages = process_data(file)

    training, testing = train_test_split(messages, test_size=0.33, shuffle=True)

    train_container = TestContainer(training)
    test_container = TestContainer(testing)

    train_x = train_container.get_messages()  # the data we give to your model
    train_y = train_container.get_is_ivan()  # the labels we want to predict

    test_x = test_container.get_messages()  # the data we give to your model
    test_y = test_container.get_is_ivan()  # the labels we want to predict

    vectorizer = CountVectorizer()

    train_x_vector = vectorizer.fit_transform(train_x)
    test_x_vector = vectorizer.transform(test_x)

    # svm = SVC(kernel='linear')
    svm = SVC(kernel="rbf")
    # svm, test_x_vector, test_y, train_x_vector, train_y = train()

    # load model from file only if the model exits
    if load_model("model.pkl"):
        svm = load_model("model.pkl")

    svm.fit(train_x_vector, train_y)

    if len(sys.argv) > 1:
        user_input = vectorizer.transform([sys.argv[1]])
        print(svm.predict(user_input))
        return

    # print(svm.predict(test_x_vector[9]))

    print(
    f1_score(test_y, svm.predict(test_x_vector), average=None, labels=[True, False])
    )
    save_model(svm, "model.pkl")


if __name__ == "__main__":
    main()
