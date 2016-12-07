# By Jack Schefer, began 12/6/16
#
from twitter import Api
from nltk.corpus import wordnet
#
##############################################################
#
# Initialization
#
with open('keys.txt') as f:
    #
    lines = f.read().splitlines()
    api = Api(  consumer_key    = lines[0].rstrip(),
            consumer_secret     = lines[1].rstrip(),
            access_token_key    = lines[2].rstrip(),
            access_token_secret = lines[3].rstrip() )
    #
#
print api.VerifyCredentials()
print
#
#
##############################################################
#
def main():
    #
    orig = input('Message to lengthen and tweet: ')
    print orig
#    to_tweet = 'Original: ' + orig + '\nLengthened:'
#    for word in orig.split():
#        #
#        to_tweet += ' ' + longest_synonym(word)
#        #
#    #
#    print to_tweet
#
#
##############################################################
#
def longest_synonym(word):
    #
    synsets = wordnet.synsets(word)
    if len(synsets == 0): return word
    syn = synsets[0].lemma_names()
    m = syn[0]
    for i in syn:
        if len(i) > len(m): m = i
    #
#
##############################################################
#
#
if __name__ == '__main__':
    #
    main()
    #
#
# End of file.
